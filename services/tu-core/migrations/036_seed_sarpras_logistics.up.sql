INSERT INTO sarpras_logistics (id, kode, nama, kategori, satuan, stok, stok_minimum) VALUES
(gen_random_uuid(), 'LOG001', 'Kertas HVS A4 80gr', 'ATK', 'Rim', 8, 10),
(gen_random_uuid(), 'LOG002', 'Spidol Whiteboard Hitam Snowman', 'ATK', 'Pcs', 45, 10),
(gen_random_uuid(), 'LOG003', 'Buku Agenda Kelas Yayasan', 'ATK', 'Buku', 3, 5),
(gen_random_uuid(), 'LOG004', 'Sabun Cuci Tangan Cair (isi Ulang)', 'Kebersihan', 'Pouch', 15, 5),
(gen_random_uuid(), 'LOG005', 'Tinta Printer Epson Hitam 003', 'Tinta', 'Botol', 1, 5)
ON CONFLICT (kode) DO NOTHING;
