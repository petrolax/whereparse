# whereparse
Простой парсер SQL-оператора 'WHERE'

Парсер работает со следующими видами строк:
1. PostgreSQL-запрос
```sql
  SELECT Alice, Bob WHERE Alice.IsActive AND Bob.LastHash = 'ab5534b'
```
2. Строка поле-значение
```sql
  Field1 = \"foo\" AND Field2 != 7 OR Field3 > 11.7
```

Парсер не поддерживает строки следующего вида:
``` sql
  SELECT Alice, Bob WHERE (Alice.IsActive OR Alice.IsDoingSomething) AND Bob.LastHash = 'ab5534b'
```
или
``` sql
  ((Field1 = \"foo\" AND Field2 != 7) OR Field3 > 11.7)
```