INSERT INTO procurement_requests (request_number, requester_id, requester_name, item_name, specification, quantity, unit, estimated_price, total_estimated, priority, justification, status)
VALUES
('REQ/2026/001', 'user-1', 'Pak Hermawan, S.Pd', 'Spidol Whiteboard', 'Snowman Hitam', 3, 'Pcs', 15000, 45000, 'NORMAL', 'Untuk mengajar Kelas 10-A', 'PENDING'),
('REQ/2026/002', 'user-2', 'Bu Linda, M.Si', 'Kertas HVS', 'A4 80gr', 2, 'Rim', 55000, 110000, 'HIGH', 'Persiapan Ujian Sekolah', 'APPROVED'),
('REQ/2026/003', 'user-3', 'Pak Danang, S.T', 'Isi Staples', 'Kenko No.10', 10, 'Kotak', 5000, 50000, 'LOW', 'Kebutuhan Administrasi Kelas', 'REJECTED')
ON CONFLICT (request_number) DO NOTHING;
