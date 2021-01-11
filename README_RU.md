
#### Общая инструкция для unix-подобных систем
1. Откройте терминал, и проверьте, установлен ли у Вас Git, выполнив команду
    ``` git version ```

    если в ответ вы получите что-то похожее
    ``` git version 2.24.3 ```
    то, можно продолжать дальше.

    Если в ответ Вы получите что другое рекомендуем вам перейти по ссылке на официальный сайт Git и следуйте официальным инструкциям по установке
    <https://git-scm.com/downloads>

2. Проверить, установлен ли у Вас компилятор Go Lang, можно выполнив команду ``` go version ```

     **Обратите внимание на версию. Актуальная версия GO не ниже 1.14!**
     
    ``` go version go1.15.2 ...```  

    Если в ответ Вы получите что другое рекомендуем вам перейти по ссылке на официальный сайт Go lang и следуйте официальным инструкциям по установке
    <https://golang.org/dl/>.

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