
#### Общая инструкция для unix-подобных систем
1. В терминале проверьте установлен ли у Вас Git. Выполните:
    ``` git version ```

    Правильный ответ будет похож на этот:
    ``` git version 2.24... ```

    Если в ответ Вы получите что-то другое, то следуйте инструкциям по установке на официальном сайте Git <https://git-scm.com/downloads>

2. Также необходим компилятор Go Lang. Выполните: ``` go version ```

     **Обратите внимание на версию. Актуальная версия GO не ниже 1.14!**
     
    ``` go version go1.15.2 ...```  

    Если в ответ Вы получите что-то другое, то следуйте инструкциям по установке на официальном сайте Go lang <https://golang.org/dl/>.

3. Клонируйте репозиторий командой ``` git clone https://github.com/tada-team/tdclient.git ```
4. Перейдите в папку проекта
     ``` cd tdclient  ```

5. Получение токена доступа

    5a. Получение токена по смс.
    ```go run examples/passwordauth/main.go -server https://demo.tada.team```
    
    5б. Получение токена по логину и паролю.
    ```go run examples/passwordauth/main.go -server https://demo.tada.team```
    
    **Если все прошло успешно получите токен**
    
    ```Your token: 7MzKGqxyqzlrxGTrozuvqEt6Qpqri26OkIApP11```

6. Далее с помощью браузера перейдите в интересующую группу. Схема URL в адресной строке браузера:

    Пример ``` https://demo.tada.team/dbd248d7-25c2-4e8f-a23a-99baf63223e9/chats/g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 ```

    Где:
    1. dbd248d7-25c2-4e8f-a23a-99baf63223e9 - идентификатор **команды** `-team`
    2. g-dce6f5fd-b741-40a6-aa9c-c0e928d9dac5 - идентификатор **чата** `-chat`
