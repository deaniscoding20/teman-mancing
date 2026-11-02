"use client";

import { useEffect, useState, useRef } from "react";
import Link from "next/link";
import { MessageSquare, Users, Send, User, Calendar } from "lucide-react";
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

export default function ChatPage() {
  const [conversations, setConversations] = useState<ConversationItem[]>([]);
  const [activeConv, setActiveConv] = useState<ConversationItem | null>(null);
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [inputText, setInputText] = useState("");
  const [loadingConvs, setLoadingConvs] = useState(true);
  const [loadingMsgs, setLoadingMsgs] = useState(false);
  const [sending, setSending] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const currentUser = getStoredUser();
  const token = getToken();

  // Load user conversations
  useEffect(() => {
    if (!token) return;
    setLoadingConvs(true);

    fetchApi<ConversationItem[]>("/conversations", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (res.data) {
          setConversations(res.data);
          if (res.data.length > 0) {
            setActiveConv(res.data[0]);
          }
        }
      })
      .catch((err) => {
        console.error("Failed to load conversations:", err);
      })
      .finally(() => {
        setLoadingConvs(false);
      });
  }, [token]);

  // Load messages when active conversation changes
  useEffect(() => {
    if (!token || !activeConv) return;
    setLoadingMsgs(true);

    fetchApi<MessageItem[]>(`/conversations/${activeConv.id}/messages`, {
      headers: { Authorization: `Bearer ${token}` },
    })
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

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token || !activeConv || !inputText.trim() || sending) return;

    setSending(true);
    try {
      const res = await fetchApi<MessageItem>(`/conversations/${activeConv.id}/messages`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ message: inputText }),
      });

      if (res.data) {
        setMessages((prev) => [...prev, res.data!]);
        setInputText("");
      }
    } catch (err: any) {
      alert(err.message || "Gagal mengirim pesan.");
    } finally {
      setSending(false);
    }
  };

  if (!currentUser) {
    return (
      <div className="max-w-md mx-auto px-4 py-20 text-center">
        <h2 className="text-xl font-bold text-brand-primary mb-2">Masuk untuk Mengakses Chat</h2>
        <p className="text-sm text-brand-muted mb-6">
          Komunikasi personal dan grup chat event hanya dapat diakses setelah login.
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-brand-surface rounded-3xl border border-brand-border h-[78vh] flex overflow-hidden shadow-sm">
        {/* Sidebar */}
        <div className="w-80 sm:w-96 border-r border-brand-border flex flex-col bg-brand-surface shrink-0">
          <div className="p-4 border-b border-brand-border flex items-center gap-2">
            <MessageSquare className="w-5 h-5 text-brand-secondary" />
            <h2 className="font-bold text-base text-brand-primary">Kotak Obrolan</h2>
          </div>

          <div className="flex-1 overflow-y-auto divide-y divide-brand-border/60">
            {loadingConvs ? (
              <div className="p-8 text-center text-xs text-brand-muted">Memuat obrolan...</div>
            ) : conversations.length === 0 ? (
              <div className="p-8 text-center text-xs text-brand-muted">
                Belum ada percakapan aktif. Ikuti janjian mancing untuk otomatis terhubung ke grup obrolan.
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
              <p className="text-sm">Pilih obrolan di sebelah kiri untuk melihat pesan</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
