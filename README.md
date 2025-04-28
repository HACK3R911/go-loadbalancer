Проверка портов на доступность

Для Windows
```
netstat -ano | findstr :9001
netstat -ano | findstr :9002
netstat -ano | findstr :9003
```

Для Linux
```
netstat -ano | findstr ":9001 :9002 :9003"
```

Для создания тестовых бекендов
```
python3 -m http.server {port}
```

Для теста нагрузки
Apache Benchmark
```
ab -n 100 -c 10 http://localhost:8080/
```