# Tugas Besar 2 IF2211 Strategi Algoritma Semester 2 2024/2025

# Oleh Kelompok: GoNext Chem

## Penjelasan Algoritma

### BFS (_Breadth First Search_)
BFS adalah algoritma pencarian yang menjelajahi graf atau pohon secara melebar, yaitu dengan mengunjungi semua simpul pada level yang sama sebelum melanjutkan ke level (kedalaman) berikutnya. Algoritma ini menggunakan struktur data queue (antrian _node_ yang dikunjungi).

### DFS (_Depth First Search_)
DFS adalah algoritma pencarian yang menjelajahi graf atau pohon dengan menyusuri cabang sedalam mungkin sebelum kembali (_backtracking_) dan menjelajahi cabang lainnya. DFS menggunakan struktur data _stack_ (tumpukan rekursi).

## _Requirement_ Program
- **Backend**: Go (Golang) dengan versi minimal 1.24.3
- **Frontend**: Node.js dan npm (Node Package Manager) dengan versi _recommended_ 10.9.2

## Instalasi dan Cara Run Aplikasi

### Backend
1. Pastikan Golang sudah diinstal.
2. Buka terminal yang mengarah pada folder _root_ (dasar) repository.
3. Pada _root_ repository, ketik 
    ```bash
    cd src/backend
    ``` 
     untuk mengakses folder backend.
4. Run command ini untuk mendownload _dependencies_:
   ```bash
   go mod tidy
   ```
5. Untuk menjalankan program backend, ketik command:
   ```bash
   go run .
   ```
6. Lanjutkan dengan menginstal dan menjalankan komponen Frontend.

### Frontend
1. Pastikan Node.js dan npm sudah diinstal.
2. Buka terminal (yang berbeda dari backend) yang mengarah pada folder _root_ (dasar) repository.
3. Pada _root_ repository, ketik 
    ```bash
    cd src/frontend
    ``` 
4. Run command ini untuk menginstal _dependencies_:
   ```bash
   npm install
   ```
5. Untuk menjalankan aplikasi React.js, ketik command:
   ```bash
   npm run start
   ```
6. Aplikasi akan membuka http://localhost:3000 secara otomatis di _browser_ default komputer. Jika tidak terbuka otomatis, klik _link_ yang tertera pada poin ini.

## Author
- Jovandra Otniel P S (13523141)
- Muhammad Aulia Azka (13523137)
- Rendi Adinata (10123083)
