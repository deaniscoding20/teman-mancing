"use client";

import { useEffect, useState, useRef, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { MessageSquare, Users, Send, User, Calendar, Plus, X, Search, Sparkles } from "lucide-react";
import { fetchApi } from "@/lib/api-client";
import { getStoredUser, getToken } from "@/lib/auth";

interface ConversationItem {
  id: string;
  type: "direct" | "group";
  event_id?: string;
  event_name?: string;
  last_message?: {
    id: string;
    message: string;
    created_at: string;
  };
}

interface MessageItem {
  id: string;
  conversation_id: string;
  sender_id: string;
  sender_name?: string;
  message: string;
  created_at: string;
}

function ChatContent() {
  const searchParams = useSearchParams();
  const paramConvId = searchParams.get("convId");
  const paramEventId = searchParams.get("eventId");
  const paramUserId = searchParams.get("userId");

  const [currentUser, setCurrentUser] = useState<any | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  const [conversations, setConversations] = useState<ConversationItem[]>([]);
  const [activeConv, setActiveConv] = useState<ConversationItem | null>(null);
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [inputText, setInputText] = useState("");
  const [loadingConvs, setLoadingConvs] = useState(true);
  const [loadingMsgs, setLoadingMsgs] = useState(false);
  const [sending, setSending] = useState(false);

  // New Chat Modal state
  const [showNewChatModal, setShowNewChatModal] = useState(false);
  const [searchUsername, setSearchUsername] = useState("");
  const [creatingChat, setCreatingChat] = useState(false);
  const [newChatError, setNewChatError] = useState<string | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setMounted(true);
    setCurrentUser(getStoredUser());
    setToken(getToken());
  }, []);

  // Load user conversations
  useEffect(() => {
    const activeToken = token || getToken();
    if (!activeToken) return;

    setLoadingConvs(true);
    fetchApi<ConversationItem[]>("/conversations")
      .then(async (res) => {
        let list = res.data || [];
        setConversations(list);

        // Check if query params specify a conversation
        if (paramConvId) {
          const found = list.find((c) => c.id === paramConvId);
          if (found) {
            setActiveConv(found);
            return;
          }
        }

        if (paramEventId) {
          let found = list.find((c) => c.event_id === paramEventId);
          if (!found) {
            try {
              const createGroupRes = await fetchApi<ConversationItem>("/conversations", {
                method: "POST",
                body: JSON.stringify({ type: "group", event_id: paramEventId }),
              });
              if (createGroupRes.data) {
                found = createGroupRes.data;
                if (!list.some((c) => c.id === found!.id)) {
                  list = [found, ...list];
                  setConversations(list);
                }
              }
            } catch (err: any) {
              console.warn("Could not auto-open group conversation:", err);
            }
          }
          if (found) {
            setActiveConv(found);
            return;
          }
        }

        if (paramUserId) {
          try {
            const createRes = await fetchApi<ConversationItem>("/conversations", {
              method: "POST",
              body: JSON.stringify({ type: "direct", participant_id: paramUserId }),
            });
            if (createRes.data) {
              const newConv = createRes.data;
              if (!list.some((c) => c.id === newConv.id)) {
                list = [newConv, ...list];
                setConversations(list);
              }
              setActiveConv(newConv);
              return;
            }
          } catch (err) {
            console.warn("Could not auto-open direct conversation:", err);
          }
        }

        if (list.length > 0 && !activeConv) {
          setActiveConv(list[0]);
        }
      })
      .catch((err) => {
        console.error("Failed to load conversations:", err);
      })
      .finally(() => {
        setLoadingConvs(false);
      });
  }, [token, paramConvId, paramEventId, paramUserId]);

  // Load messages when active conversation changes
  useEffect(() => {
    const activeToken = token || getToken();
    if (!activeToken || !activeConv) return;

    setLoadingMsgs(true);
    fetchApi<MessageItem[]>(`/conversations/${activeConv.id}/messages`)
      .then((res) => {
        if (res.data) setMessages(res.data);
      })
      .catch((err) => {
        console.error("Failed to load messages:", err);
      })
      .finally(() => {
        setLoadingMsgs(false);
      });
  }, [activeConv, token]);

  // Scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Real-time WebSocket connection
  useEffect(() => {
    const activeToken = token || getToken();
    if (!activeToken) return;

    const baseWsUrl = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1/ws";
    const wsUrl = `${baseWsUrl}?token=${encodeURIComponent(activeToken)}`;

    let ws: WebSocket | null = null;
    try {
      ws = new WebSocket(wsUrl);
      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload && (payload.type === "new_message" || payload.type === "chat_message")) {
            const incoming: MessageItem = {
              id: payload.data?.id || payload.message_id || payload.id || `ws-${Date.now()}`,
              conversation_id: payload.data?.conversation_id || payload.conversation_id,
              sender_id: payload.data?.sender_id || payload.sender_id,
              sender_name: payload.data?.sender_name || payload.sender_name,
              message: payload.data?.message || payload.content || payload.message || "",
              created_at: payload.data?.created_at || payload.created_at || new Date().toISOString(),
            };

            if (activeConv && incoming.conversation_id === activeConv.id) {
              setMessages((prev) => {
                if (prev.some((m) => m.id === incoming.id)) return prev;
                return [...prev, incoming];
              });
            }
          }
        } catch (e) {
          // ignore non-json messages
        }
      };
    } catch (err) {
      console.warn("WebSocket connection not available:", err);
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [token, activeConv]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    const activeToken = token || getToken();
    if (!activeToken || !activeConv || !inputText.trim() || sending) return;

    if (currentUser && !currentUser.is_verified) {
      alert("Akun Anda belum diverifikasi. Silakan periksa email Anda untuk verifikasi sebelum mengirim pesan.");
      return;
    }

    const textToSend = inputText.trim();
    setSending(true);

    try {
      const res = await fetchApi<MessageItem>(`/conversations/${activeConv.id}/messages`, {
        method: "POST",
        body: JSON.stringify({ message: textToSend }),
      });

      if (res.data) {
        const newMsg = res.data;
        setMessages((prev) => {
          if (prev.some((m) => m.id === newMsg.id)) return prev;
          return [...prev, newMsg];
        });
        setInputText("");

        // Update last_message in conversation list
        setConversations((prev) =>
          prev.map((c) =>
            c.id === activeConv.id
              ? {
                  ...c,
                  last_message: {
                    id: newMsg.id,
                    message: newMsg.message,
                    created_at: newMsg.created_at,
                  },
                }
              : c
          )
        );
      }
    } catch (err: any) {
      alert(err.message || "Gagal mengirim pesan.");
    } finally {
      setSending(false);
    }
  };

  const handleStartDirectChat = async (e: React.FormEvent) => {
    e.preventDefault();
    const activeToken = token || getToken();
    if (!activeToken || !searchUsername.trim() || creatingChat) return;

    setCreatingChat(true);
    setNewChatError(null);

    try {
      const userRes = await fetchApi<{ id: string; username: string; name: string }>(
        `/users/${searchUsername.trim().toLowerCase()}`
      );
      if (!userRes.data?.id) {
        throw new Error("Pengguna dengan username tersebut tidak ditemukan.");
      }

      const targetUserId = userRes.data.id;
      if (currentUser && currentUser.id === targetUserId) {
        throw new Error("Anda tidak dapat memulai obrolan dengan akun sendiri.");
      }

      const convRes = await fetchApi<ConversationItem>("/conversations", {
        method: "POST",
        body: JSON.stringify({
          type: "direct",
          participant_id: targetUserId,
        }),
      });

      if (convRes.data) {
        const conv = convRes.data;
        setConversations((prev) => {
          if (prev.some((c) => c.id === conv.id)) return prev;
          return [conv, ...prev];
        });
        setActiveConv(conv);
        setShowNewChatModal(false);
        setSearchUsername("");
      }
    } catch (err: any) {
      setNewChatError(err.message || "Gagal memulai obrolan baru.");
    } finally {
      setCreatingChat(false);
    }
  };

  if (!mounted) {
    return <div className="text-center py-20 text-brand-muted">Memuat kotak obrolan...</div>;
  }

  if (!currentUser) {
    return (
      <div className="max-w-md mx-auto px-4 py-20 text-center">
        <h2 className="text-xl font-bold text-brand-primary mb-2">Masuk untuk Mengakses Chat</h2>
        <p className="text-sm text-brand-muted mb-6">
          Komunikasi personal dan grup chat event hanya dapat diakses setelah login.
        </p>
        <Link
          href="/login"
          className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-brand-secondary text-white font-semibold text-sm shadow-sm hover:bg-brand-secondary/90 transition"
        >
          Masuk ke Akun
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-brand-surface rounded-3xl border border-brand-border h-[78vh] flex overflow-hidden shadow-sm">
        {/* Sidebar */}
        <div className="w-80 sm:w-96 border-r border-brand-border flex flex-col bg-brand-surface shrink-0">
          <div className="p-4 border-b border-brand-border flex items-center justify-between">
            <div className="flex items-center gap-2">
              <MessageSquare className="w-5 h-5 text-brand-secondary" />
              <h2 className="font-bold text-base text-brand-primary">Kotak Obrolan</h2>
            </div>
            <button
              onClick={() => setShowNewChatModal(true)}
              className="p-1.5 rounded-xl border border-brand-border bg-brand-background hover:bg-brand-secondary hover:text-white text-brand-primary transition text-xs font-semibold flex items-center gap-1 shadow-xs"
              title="Mulai Chat Baru"
            >
              <Plus className="w-4 h-4" />
              <span className="hidden sm:inline">Chat Baru</span>
            </button>
          </div>

          <div className="flex-1 overflow-y-auto divide-y divide-brand-border/60">
            {loadingConvs ? (
              <div className="p-8 text-center text-xs text-brand-muted">Memuat obrolan...</div>
            ) : conversations.length === 0 ? (
              <div className="p-8 text-center flex flex-col items-center">
                <MessageSquare className="w-10 h-10 text-brand-muted/40 mb-3" />
                <p className="text-xs text-brand-muted mb-4">
                  Belum ada percakapan aktif. Mulai obrolan langsung atau ikuti janjian mancing.
                </p>
                <button
                  onClick={() => setShowNewChatModal(true)}
                  className="px-4 py-2 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 transition shadow-sm"
                >
                  + Mulai Chat Baru
                </button>
              </div>
            ) : (
              conversations.map((c) => {
                const isActive = activeConv?.id === c.id;
                const title = c.type === "group" ? (c.event_name || "Grup Event") : "Obrolan Personal";

                return (
                  <button
                    key={c.id}
                    onClick={() => setActiveConv(c)}
                    className={`w-full text-left p-4 flex items-start gap-3 hover:bg-brand-background transition ${
                      isActive ? "bg-brand-background border-l-4 border-l-brand-secondary" : ""
                    }`}
                  >
                    <div className="w-10 h-10 rounded-xl bg-brand-primary/10 flex items-center justify-center text-brand-secondary shrink-0">
                      {c.type === "group" ? (
                        <Users className="w-5 h-5" />
                      ) : (
                        <User className="w-5 h-5" />
                      )}
                    </div>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between gap-1 mb-1">
                        <span className="font-bold text-xs text-brand-primary truncate">
                          {title}
                        </span>
                        <span className="text-[10px] text-brand-muted uppercase font-semibold">
                          {c.type}
                        </span>
                      </div>
                      <p className="text-xs text-brand-muted truncate">
                        {c.last_message?.message || "Belum ada pesan"}
                      </p>
                    </div>
                  </button>
                );
              })
            )}
          </div>
        </div>

        {/* Chat Window */}
        <div className="flex-1 flex flex-col bg-brand-background">
          {activeConv ? (
            <>
              {/* Header */}
              <div className="p-4 bg-brand-surface border-b border-brand-border flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-brand-primary/10 flex items-center justify-center text-brand-secondary">
                    {activeConv.type === "group" ? <Users className="w-5 h-5" /> : <User className="w-5 h-5" />}
                  </div>
                  <div>
                    <h3 className="font-bold text-sm text-brand-primary">
                      {activeConv.type === "group" ? (activeConv.event_name || "Grup Event Mancing") : "Obrolan Personal"}
                    </h3>
                    <p className="text-[10px] text-brand-muted">
                      {activeConv.type === "group" ? "Grup Terbuka Anggota Event" : "Pesan Terenkripsi"}
                    </p>
                  </div>
                </div>

                {activeConv.event_id && (
                  <Link
                    href={`/events/${activeConv.event_id}`}
                    className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-brand-border bg-brand-background text-brand-primary text-xs font-semibold hover:border-brand-secondary transition"
                  >
                    <Calendar className="w-3.5 h-3.5 text-brand-secondary" />
                    Detail Event
                  </Link>
                )}
              </div>

              {/* Message History */}
              <div className="flex-1 p-4 sm:p-6 overflow-y-auto space-y-4">
                {loadingMsgs ? (
                  <div className="text-center py-10 text-xs text-brand-muted">Memuat pesan...</div>
                ) : messages.length === 0 ? (
                  <div className="text-center py-16 text-xs text-brand-muted">
                    Mulai percakapan dengan menyapa rekan pemancingmu di sini!
                  </div>
                ) : (
                  messages.map((m) => {
                    const isMe = m.sender_id === currentUser.id;
                    return (
                      <div
                        key={m.id}
                        className={`flex flex-col ${isMe ? "items-end" : "items-start"}`}
                      >
                        <div className="flex items-center gap-1.5 mb-1 px-1">
                          <span className="text-[10px] font-bold text-brand-muted">
                            {isMe ? "Anda" : (m.sender_name || "Pemancing")}
                          </span>
                          <span className="text-[10px] text-brand-muted/70">
                            {new Date(m.created_at).toLocaleTimeString("id-ID", {
                              hour: "2-digit",
                              minute: "2-digit",
                            })}
                          </span>
                        </div>
                        <div
                          className={`max-w-md sm:max-w-lg px-4 py-2.5 rounded-2xl text-xs sm:text-sm leading-relaxed ${
                            isMe
                              ? "bg-brand-primary text-white rounded-br-xs shadow-sm"
                              : "bg-brand-surface text-brand-text border border-brand-border rounded-bl-xs shadow-xs"
                          }`}
                        >
                          {m.message}
                        </div>
                      </div>
                    );
                  })
                )}
                <div ref={messagesEndRef} />
              </div>

              {/* Input Area */}
              <form onSubmit={handleSendMessage} className="p-4 bg-brand-surface border-t border-brand-border flex items-center gap-2">
                <input
                  type="text"
                  required
                  value={inputText}
                  onChange={(e) => setInputText(e.target.value)}
                  placeholder="Ketik pesan untuk dikirimkan..."
                  className="flex-1 px-4 py-2.5 rounded-xl border border-brand-border bg-brand-background text-sm text-brand-text placeholder-brand-muted/70 focus:outline-none focus:ring-2 focus:ring-brand-secondary transition"
                />
                <button
                  type="submit"
                  disabled={sending || !inputText.trim()}
                  className="p-2.5 rounded-xl bg-brand-secondary text-white hover:bg-brand-secondary/90 disabled:opacity-50 transition shadow-sm"
                  title="Kirim Pesan"
                >
                  <Send className="w-5 h-5" />
                </button>
              </form>
            </>
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-center p-8 text-brand-muted">
              <MessageSquare className="w-12 h-12 text-brand-muted/40 mb-3" />
              <p className="text-sm font-semibold mb-2">Pilih obrolan di sebelah kiri untuk melihat pesan</p>
              <button
                onClick={() => setShowNewChatModal(true)}
                className="mt-2 px-4 py-2 rounded-xl border border-brand-border bg-brand-surface text-xs font-semibold text-brand-primary hover:border-brand-secondary transition"
              >
                Atau Cari Angler untuk Chat Baru
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Modal Mulai Chat Baru */}
      {showNewChatModal && (
        <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-brand-surface rounded-3xl border border-brand-border max-w-md w-full p-6 shadow-xl animate-in fade-in zoom-in-95 duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <Sparkles className="w-5 h-5 text-brand-secondary" />
                <h3 className="text-base font-bold text-brand-primary">Mulai Chat Personal Baru</h3>
              </div>
              <button
                onClick={() => {
                  setShowNewChatModal(false);
                  setNewChatError(null);
                }}
                className="p-1 rounded-lg text-brand-muted hover:text-brand-primary hover:bg-brand-background transition"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <p className="text-xs text-brand-muted mb-4">
              Ketik username angler yang ingin Anda ajak mengobrol (misal: <code>budi_angler</code> atau <code>deanmancing</code>):
            </p>

            <form onSubmit={handleStartDirectChat} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-brand-text mb-1">Username Angler</label>
                <div className="relative">
                  <span className="absolute left-3 top-2.5 text-brand-muted text-sm font-bold">@</span>
                  <input
                    type="text"
                    required
                    value={searchUsername}
                    onChange={(e) => setSearchUsername(e.target.value)}
                    placeholder="username"
                    className="w-full pl-8 pr-4 py-2 rounded-xl border border-brand-border bg-brand-background text-sm text-brand-text focus:outline-none focus:ring-2 focus:ring-brand-secondary"
                  />
                </div>
              </div>

              {newChatError && (
                <div className="p-3 rounded-xl bg-red-50 text-red-600 text-xs border border-red-200">
                  {newChatError}
                </div>
              )}

              <div className="flex justify-end gap-2 pt-2">
                <button
                  type="button"
                  onClick={() => setShowNewChatModal(false)}
                  className="px-4 py-2 rounded-xl border border-brand-border text-xs font-semibold text-brand-muted hover:text-brand-primary transition"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={creatingChat || !searchUsername.trim()}
                  className="px-5 py-2 rounded-xl bg-brand-secondary text-white text-xs font-bold hover:bg-brand-secondary/90 disabled:opacity-50 transition shadow-sm"
                >
                  {creatingChat ? "Mencari..." : "Mulai Obrolan"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default function ChatPage() {
  return (
    <Suspense fallback={<div className="text-center py-20 text-brand-muted">Memuat kotak obrolan...</div>}>
      <ChatContent />
    </Suspense>
  );
}
