# RR http-балансировщик

Для начала работы нужно:

1. Проверить порты на доступность

**Для Windows**
```powershell
netstat -ano | findstr :8080
netstat -ano | findstr :9001
netstat -ano | findstr :9002
netstat -ano | findstr :9003
```

**Для Linux**
```bash
netstat -ano | findstr ":8080 :9001 :9002 :9003"
```

Для создания тестовых бекендов 
```
python -m http.server {port}
```

Для теста нагрузки
Apache Benchmark
```
ab -n 100 -c 10 http://localhost:8080/
```