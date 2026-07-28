<template>
  <div class="fullscreen bg-blue-grey-1 column flex-center q-pa-md">
    <div class="text-center column items-center q-gutter-md" style="max-width: 400px;">
      <!-- Ikon Peringatan -->
      <q-icon name="gpp_bad" size="120px" color="negative" />

      <!-- Pesan Error -->
      <div class="text-h4 text-weight-bold text-grey-9">Akses Ditolak</div>
      <div class="text-body1 text-grey-7">
        Sesi Anda telah berakhir atau Anda tidak memiliki izin untuk mengakses halaman ini. Silakan masuk kembali.
      </div>

      <!-- Tombol Kembali ke Beranda -->
      <q-btn unelevated no-caps color="primary" label="Kembali ke Beranda" to="/"
        class="q-px-lg q-py-sm text-weight-bold q-mt-md" />
    </div>
  </div>
</template>

