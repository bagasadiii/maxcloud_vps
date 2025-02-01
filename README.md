# Cloudmax_VPS
Ini adalah project untuk test backend ke perusahaan Maxcloud, aplikasi ini digunakan untuk monitoring VPS client yang masuk dan otomatis memotong saldo setiap satu jam.
Untuk menjalankan aplikasi, clone aplikasi ini menggunakan
```sh
git clone https://github.com/bagasadiii/maxcloud_vps
cd maxcloud_vps
```
Lalu jalankan docker dengan perintah
```sh
docker compose up -d --build
```
Aplikasi akan otomatis berjalan melalui http://localhost:8080

## Mendaftarkan client
Untuk menyimpan client, saya menggunakan database PostgreSQL untuk menyimpan data client.
Anda bisa mendaftarkan client dengan cara mengirim request JSON ke endpoint http://localhost:8080/api/register dengan method post.
ada tiga komponen yang dikirimkan ke JSON yaitu **Email, Plan dan Balance**
Untuk value plan saat ini hanya bisa menggunakan plan basic, normal dan premium
Untuk value balance mohon masukkan angka yang lebih dari down payment, karena balance diperlukan untuk membayar down payment di awal service

## Contoh
```json
{
  "email": "test@example.com",
  "balance": 3000000,
  "plan": "basic"
}
```

Ketika berhasil di daftarkan, aplikasi akan memonitor saldo setiap jam dan melakukan transaksi pemotongan saldo secara otomatis tergantung Cost Per Hour yang sudah ditentukan di plan

