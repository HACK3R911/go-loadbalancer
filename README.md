# RR http-балансировщик

[Тестовое задание Cloudru](https://github.com/Go-Cloud-Camp/test-assignment/tree/main)

[Ответы на вопросы](_question1-3.md)  

[Функционал](#Функционал)  

[Начало работы](#Начало-работы)  

## Функционал
**Основной функционал:**
- Распределение запросов (Round-Robin)
- Параллелизм и конкурентность (горутины каналы)
- Обработка ошибок
- Логирование (log)
- Конфигурация (config.yaml)

**Что реализовано, помимо основго функционала:**
- Health Checks
- Graceful Shutdown

**Что использовал внешнего:**
- Для парсинга конфиг файла
    gopkg.in/yaml.v3 

## Начало работы

**1. Проверить порты на доступность**

Для Windows
```powershell
netstat -ano | findstr :8080
netstat -ano | findstr :9001
netstat -ano | findstr :9002
netstat -ano | findstr :9003
```

Для Linux
```bash
netstat -ano | findstr ":8080 :9001 :9002 :9003"
```

Если все порты свободные переходим на следующий шаг, иначе можно поменять порты в config.yaml файле 

**2. Для развертывания через Docker**

```
docker-compose up --build
```

**3. Теста работы балансировщика**
Установка Apache Benchmark через scoop для Windows:
```
scoop install httpd
```

Для Debian-совместимой ОС:
```
sudo apt install apache2-utils
```

Apache Benchmark
```
ab -n 100 -c 10 http://localhost:8080/ # Если меняли конфиг, то изменяем порт
```