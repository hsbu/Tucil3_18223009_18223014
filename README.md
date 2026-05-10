# Tugas Kecil 3 IF2211 Strategi Algoritma Ice Sliding Puzzle Solver

## Author

- 18223009 Muhammad Faiz Alfikrona
- 18223014 Muhamad Hasbullah Faris

## Deskripsi Singkat

Ice Sliding Puzzle adalah permainan logika di mana pemain harus menggerakkan karakter dari titik awal menuju titik keluar di atas permukaan es yang licin. Pin hanya dapat bergerak secara horizontal atau vertikal, namun karena permukaan licin, karakter tidak akan berhenti bergerak sampai menabrak dinding atau rintangan.

Program ini menyelesaikan puzzle ice sliding dengan checkpoint berurutan. Pemain mulai dari `Z`, harus mengunjungi checkpoint `0`, `1`, `2`, dan seterusnya secara urut, lalu mencapai goal `O`. Setiap langkah berupa slide sampai berhenti di dinding, goal, atau kondisi tidak valid seperti lava `L`.

Algoritma yang tersedia:

- UCS, yang dikenal juga sebagai Djikstra
- Greedy Best-First Search, hanya memperhitungkan heuristic
- A*, memperhitungkan cost existing dan heuristic
- IDA*, menjalaknkan A* dengan pendekatan iterative deepening

Program tersedia dalam mode CLI dan GUI.

## Requirement

- Go `1.25.6` atau versi kompatibel sesuai `src/go.mod`
- Dependency Go akan diunduh otomatis oleh `go run`, `go build`, atau `go mod download`
- Untuk GUI digunakan library Ebiten (`github.com/hajimehoshi/ebiten/v2`)

## Cara Kompilasi

Jalankan dari root repository:

```powershell
cd src
go build -o ../bin/cli.exe ./cmd/cli
go build -o ../bin/gui.exe ./cmd/gui
```

Hasil build akan menghasilkan executable CLI dan GUI di folder `bin`.

## Cara Menjalankan

### CLI

```powershell
cd src
go run ./cmd/cli
```

Masukkan path file input, misalnya:

```text
../test/sample5.txt
```

Setelah solusi ditemukan, program menampilkan solusi, cost, waktu eksekusi dalam millisecond, dan jumlah iterasi. Jika memilih playback, gunakan:

- `A`/`D`: mundur/maju step
- `ESC` atau `J`: lompat ke step tertentu
- `Q`: keluar playback

### GUI

```powershell
cd src
go run ./cmd/gui
```

Pada input file GUI, cukup masukkan nama file dari folder `test`, misalnya:

```text
sample5.txt
```

Alur penggunaan GUI:

1. Masukkan file input.
2. Pilih algoritma.
3. Jika algoritma membutuhkan heuristic, pilih heuristic.
4. Lihat hasil solusi.
5. Gunakan `Playback >>` untuk visualisasi langkah.
6. Gunakan `Save TXT` untuk menyimpan hasil ke `solution.txt`.
