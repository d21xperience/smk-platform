CREATE TABLE IF NOT EXISTS "tabel_jadwal_pelajaran" (
	"id" UUID NOT NULL,
	"pembelajaran_id" UUID NOT NULL REFERENCES "tabel_pembelajaran"("pembelajaran_id") ON DELETE CASCADE,
	"hari" VARCHAR(10) NOT NULL,
	"jam_mulai" TIME NOT NULL,
	"jam_selesai" TIME NOT NULL,
	"created_at" TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY ("id")
);
