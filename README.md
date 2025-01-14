# Работа с миграциями

Команда применяет миграции в порядке номера:

```bash
make migrate-up
```

Команда откатывает последнюю миграцию:

```bash
make migrate-down
```

Команда устанавливает текущую версию БД. Нужна для очищения состояния и восстановления порядка:

```bash
make migrate-force VERSION=n
```  

Помогает устранить ошибку: - _Dirty database version 1. Fix and force version_

Команда откатывает миграции на n-ое количество шагов назад:

```bash
make migrate-down-stepback STEPBACK=n
```

Команда откатывает все миграции:

```bash
make migrate-down-all
```
