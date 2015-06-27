# controllers

## controller.go

#### GetModel

Возвращает модель по строке - названию сущности БД.

#### Render

Исполнение шаблона из папки "mvc/views/".

#### CheckSid

Проверяет наличие в БД пользователя с SID-ом из кук. Возвращает ID пользователя и ошибку.

#### isAdmin

Возвращает истину или ложь: проверка роли пользователя.

## auth.go

#### HandleRegister

Регистрация пользователя в системе. Создание `физического лица` и `регистрации`.
Пароль пользователя хранится в базе в неявном виде. В качестве соли для генерации хеша пароля используется текущее время.

#### HandleLogin

Вход в систему. Начало сеанса: генерация SID, хранение его в куках.

#### HandleLogout

Выход из системы. Окончание сеанса: чистка SID и кук.

## form-import.go

#### GetEventTypesByEventId

Возвращает названия типов мероприятий по ID мероприятия.

#### ImportForms

Копирование форм из последнесозданных мероприятий указанных типов мероприятий.

## person-list.go

#### GetPersonsByEventId

Список значений об участниках мероприятия по параметрам анкет.

#### GetParamsByEventId

Список параметров по идентификатору мероприятия.

## index.go

#### Index

Загрузка стартовой страницы сайта. Анонс мероприятий.

#### Init

Сброс БД. Загрузка тестовых данных.

#### LoadContestsFromCats

Загрузка в БД незавершенных мероприятий из CATS.

#### CreateRegistrationEvent

Создание мероприятия, соответсвующих форм и полей для регистрации пользователей в системе.

## item.go

#### GetEditHistoryData

#### GetHistoryRequest

#### GetListHistoryEvents

#### GetRequest

#### RegPerson

#### InsertUserParams

## load.go

#### Load

Подгрузка данных в грид для любой таблицы БД.

#### UserGroupsLoad

Подгрузка данных о группах, созданных данным пользователем.

#### GroupsLoad

Подгрузка данных о группах, членом которых является данный пользователь.

#### RegistrationsLoad

Подгрузка данных о регистрациях пользователя.

#### GroupRegistrationsLoad

Подгрузка данных о регистрациях групп, зарегестрированных пользователем.

#### PersonsLoad

Подгрузка данных об участниках группы.

## blank.go

#### GetPersonRequestFromGroup

#### GetPersonRequest

#### ConfirmOrRejectPersonRequest

#### EditParams

#### AddPerson

## groups.go

#### RegGroup

#### ConfirmInvitationToGroup

#### RejectInvitationToGroup

#### IsRegGroup

## handler.go

#### GetList

#### Index

#### ShowCabinet

#### WellcomeToProfile

#### Login

#### CheckEnableOfUser

## grid-handler.go

#### GetSubTable

#### CreateGrid

#### EditGridRow

#### JsonToExcel
