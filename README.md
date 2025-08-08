
# XSS Zaafiyeti Senaryosu


## Yükleme 


Projeyi klonlayın

```bash
  git clone https://github.com/Eneskalin/XSS-Basic-VulnLab.git
```

Proje dizinine gidin

```bash
  cd XSS-Basic-VulnLab
```

user.json üzerinden kullanıcıları tanımlayın
```bash
  [
    {"username": "admin", "password": "admin123", "role": "admin"},
    {"username": "user1", "password": "user123", "role": "user"},
    {"username": "user2", "password": "user123", "role": "user"},
    {"username": "user3", "password": "user123", "role": "user"},
    {"username": "user4", "password": "user123", "role": "user"},
    {"username": "test", "password": "test", "role": "user"}
  ]
  
```
Paketleri yükleyin

```bash
  go mod tidy
```

Sunucuyu çalıştırın

```bash
  go run main.go
```

  
## Çözüm

[Link](https://eneskalin.com/blog/xss-vulnlab-solution)


  
