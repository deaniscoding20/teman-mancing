-- 000002_seed_initial_data.up.sql
-- Seed standard fish types and fishing techniques common in Indonesia

INSERT INTO fish_types (id, name) VALUES
    (gen_random_uuid(), 'Kakap Putih (Barramundi)'),
    (gen_random_uuid(), 'Kakap Merah'),
    (gen_random_uuid(), 'Kerapu'),
    (gen_random_uuid(), 'Tenggiri'),
    (gen_random_uuid(), 'Tongkol'),
    (gen_random_uuid(), 'Tuna'),
    (gen_random_uuid(), 'Giant Trevally (GT)'),
    (gen_random_uuid(), 'Barakuda'),
    (gen_random_uuid(), 'Nila'),
    (gen_random_uuid(), 'Lele'),
    (gen_random_uuid(), 'Patin'),
    (gen_random_uuid(), 'Bawal'),
    (gen_random_uuid(), 'Gabus (Channa)'),
    (gen_random_uuid(), 'Gurame'),
    (gen_random_uuid(), 'Mas')
ON CONFLICT (name) DO NOTHING;

INSERT INTO fishing_techniques (id, name) VALUES
    (gen_random_uuid(), 'Casting'),
    (gen_random_uuid(), 'Dasaran (Bottom Fishing)'),
    (gen_random_uuid(), 'Jigging'),
    (gen_random_uuid(), 'Popping'),
    (gen_random_uuid(), 'Trolling'),
    (gen_random_uuid(), 'Ngoncer (Live Baiting)'),
    (gen_random_uuid(), 'Fly Fishing'),
    (gen_random_uuid(), 'Pelampung (Float Fishing)'),
    (gen_random_uuid(), 'Surf Casting')
ON CONFLICT (name) DO NOTHING;
