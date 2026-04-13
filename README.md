---

# 🚀 mysubdo - Subdomain Automation Tool

`mysubdo` adalah tool berbasis Golang untuk melakukan **subdomain enumeration otomatis** dari berbagai sumber, kemudian menggabungkan hasilnya dan melakukan pengecekan host aktif menggunakan `httpx`.

---

## ✨ Features

* 🔍 Subdomain enumeration dari multiple sources:

  * assetfinder
  * subfinder
  * crt.sh
  * anubis
  * hackertarget
  * urlscan
* ⚡ Concurrent execution (goroutine)
* 🧹 Deduplication otomatis
* 🌐 Live host detection dengan `httpx`
* 📄 Output terstruktur

---

## 📦 Requirements

Pastikan tools berikut sudah terinstall:

```bash
assetfinder
subfinder
httpx
```

Install contoh:

```bash
go install github.com/tomnomnom/assetfinder@latest
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/projectdiscovery/httpx/cmd/httpx@latest
```

---

## ⚙️ Installation

Clone repo:

```bash
git clone https://github.com/yourusername/mysubdo.git
cd mysubdo
```

Build:

```bash
go build -o mysubdo
```

---

## 🚀 Usage

```bash
./mysubdo target.com
```

---

## 📂 Output Files

| File               | Deskripsi                            |
| ------------------ | ------------------------------------ |
| `assetfinder.txt`  | hasil dari assetfinder               |
| `subfinder.txt`    | hasil dari subfinder                 |
| `crtsh.txt`        | hasil dari crt.sh                    |
| `anubis.txt`       | hasil dari anubis                    |
| `hackertarget.txt` | hasil dari hackertarget              |
| `urlscan.txt`      | hasil dari urlscan                   |
| `allsubs.txt`      | semua subdomain (deduplicated)       |
| `httpx.txt`        | hasil live host + informasi tambahan |

---

## 🔎 Workflow

1. Input domain
2. Jalankan semua source secara parallel
3. Ambil subdomain dari:

   * API
   * CLI tools
4. Merge & deduplicate
5. Jalankan `httpx` untuk:

   * status code
   * title
   * tech detect
6. Simpan hasil

---