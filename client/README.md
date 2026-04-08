# rsl.pos.client

Go-клиент для SOAP-сервиса RSLoyaltyService (WCF, SOAP 1.2 + WS-Addressing).

## Установка

```bash
go get github.com/ramzes4rules/rsl6.pos.client
```

## Использование

### Создание клиента

```go
package main

import (
    "context"
    "crypto/tls"
    "log"
    "net/http"
    "time"

    rslpos "github.com/ramzes4rules/rsl6.pos.client"
)

func main() {
    // Простое создание
    client := rslpos.NewClient("https://server/RS.Loyalty.Service/RSLoyaltyService.svc")

    // С настройками
    client = rslpos.NewClient(
        "https://server/RS.Loyalty.Service/RSLoyaltyService.svc",
        rslpos.WithTimeout(15*time.Second),
        rslpos.WithHTTPClient(&http.Client{
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
            },
        }),
    )

    ctx := context.Background()

    // Ping
    if err := client.Ping(ctx); err != nil {
        log.Fatal(err)
    }

    // Проверка статуса
    online, err := client.IsOnline(ctx, "2.0")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Online: %v", online)

    // Получение баланса карты
    balance, err := client.GetCardBalance(ctx, "1234567890")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Balance: %s", balance)
}
```

### Обработка ошибок SOAP

```go
err := client.Ping(ctx)
if fe, ok := rslpos.IsFaultError(err); ok {
    log.Printf("SOAP Fault: code=%s reason=%s detail=%s", fe.Code, fe.Reason, fe.Detail)
}
```

### Использование mock в тестах

```go
package myservice_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    rslpos "github.com/ramzes4rules/rsl6.pos.client"
    rslposmock "github.com/ramzes4rules/rsl6.pos.client/mock"
)

func TestMyBusinessLogic(t *testing.T) {
    m := new(rslposmock.MockService)
    m.On("GetCardBalance", mock.Anything, "CARD001").Return("1500.50", nil)
    m.On("IsCardValid", mock.Anything, "CARD001").Return(true, nil)

    // Передаём mock как rslpos.Service в вашу бизнес-логику
    result := processCard(m, "CARD001")
    assert.Equal(t, "1500.50", result)
    m.AssertExpectations(t)
}

func processCard(svc rslpos.Service, card string) string {
    ctx := context.Background()
    valid, _ := svc.IsCardValid(ctx, card)
    if !valid {
        return ""
    }
    balance, _ := svc.GetCardBalance(ctx, card)
    return balance
}
```

## Архитектура (Clean Architecture)

```
rsl.pos.client/
│
│  ┌─── Domain Layer (публичный контракт, без зависимостей) ───┐
├── doc.go             # Документация пакета и обзор архитектуры
├── models.go          # Доменные типы: StoreConfig, RSInfoPacket, ItemStatistics
├── service.go         # Интерфейс Service (порт) — 36 операций
├── errors.go          # FaultError — доменная ошибка SOAP Fault
│  └──────────────────────────────────────────────────────────┘
│
│  ┌─── Client Layer (адаптер, связывает домен и инфраструктуру) ─┐
├── client.go          # Client struct, NewClient, Option, HTTP-транспорт
├── operations.go      # 36 методов Client (SOAP-операции)
│  └────────────────────────────────────────────────────────────┘
│
│  ┌─── Infrastructure Layer (скрыт через internal/) ──────────┐
├── internal/
│   └── soap/
│       ├── envelope.go   # BuildEnvelope, ParseEnvelope, RawFault
│       ├── xmltypes.go   # XmlDecimal, XmlDateTime, XmlStoreConfig*, XmlLongArray
│       └── soap_test.go  # Юнит-тесты инфраструктуры
│  └──────────────────────────────────────────────────────────┘
│
│  ┌─── Mock Layer (для тестирования потребителей) ─────────────┐
├── mock/
│   ├── mock_service.go      # MockService (testify/mock)
│   └── mock_service_test.go # Тесты mock
│  └──────────────────────────────────────────────────────────┘
│
├── client_test.go     # Интеграционные тесты клиента (httptest)
├── docs/
│   └── RSLoyaltyService.wsdl
├── go.mod
└── go.sum
```

### Принципы

- **Dependency Rule**: доменный слой (`models.go`, `service.go`, `errors.go`) не зависит ни от чего внешнего
- **Ports & Adapters**: `Service` — порт, `Client` — адаптер, `internal/soap` — инфраструктура
- **internal/**: SOAP-детали спрятаны через Go convention `internal/`, потребители не могут импортировать
- **Тестируемость**: mock реализует интерфейс `Service` — потребители тестируют бизнес-логику без SOAP

## Поддерживаемые операции (36)

| Операция | Описание |
|---|---|
| `Ping` | Проверка доступности |
| `IsOnline` | Статус онлайн |
| `GetParameters` | Получение параметров |
| `GetStoreSettings` | Настройки магазина |
| `RegisterDiscountCard` | Регистрация карты |
| `IsCardValid` | Проверка карты |
| `IsCouponValid` | Проверка купона |
| `GetVerifyCode` | Код верификации |
| `GetCardBalance` | Баланс карты |
| `GetCardDiscountAmount` | Сумма скидки (decimal) |
| `GetCardDiscountAmountString` | Сумма скидки (string) |
| `GetMessages` | Сообщения |
| `GetDiscounts` | Скидки |
| `GetEmail` | Email клиента |
| `GetSelfBuyDiscounts` | Скидки по самовыкупу |
| `Accrual` | Начисление бонусов |
| `OfflineAccrual` | Офлайн-начисление |
| `Refund` | Возврат |
| `SubtractBonus` | Списание бонусов |
| `SubtractBonus45` | Списание бонусов v4.5 |
| `CancelSubtractBonus` | Отмена списания |
| `ValidateUser` | Валидация пользователя |
| `CheckDiscountCard` | Проверка дисконтной карты |
| `ValidateUserRole` | Проверка роли |
| `GetUserRole` | Получение роли |
| `ActivationPaymentCard` | Активация платёжной карты |
| `CancelActivationPaymentCard` | Отмена активации |
| `QuerySyncStream` | Запрос потока синхронизации |
| `GetSyncStream` | Получение потока синхронизации |
| `IsTaskCompleted` | Проверка статуса задачи |
| `GetUpdateStream` | Поток обновления |
| `UploadReferences` | Загрузка справочников |
| `GetReferencesStamp` | Штамп справочников |
| `GetDataPacket` | Пакет данных |
| `SendInfoPacket` | Отправка инфо-пакета |
| `GetStatistic` | Статистика |

## Тестирование

```bash
go test ./... -v
```

