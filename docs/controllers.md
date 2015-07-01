# controllers

![controllers][scheme-1]

UPDATE-1
--------

![controllers][scheme-2]

## Controller

#### GetModel

Возвращает модель по строке - названию сущности БД.

#### Render

Исполнение шаблона из папки "mvc/views/".

#### CheckSid

Проверяет наличие в БД пользователя с SID-ом из кук. Возвращает ID пользователя и ошибку.

#### isAdmin

Возвращает истину или ложь: проверка роли пользователя.

## RegistrationController

#### EventRegisterAction

Регистрация пользователя в мероприятии.

#### InsertUserParams

Сохранение папаметров анкеты.

#### Register

Регистрация пользователя в системе. Создание `физического лица` и `регистрации`.
Пароль пользователя хранится в базе в неявном виде. В качестве соли для генерации хеша пароля используется текущее время.

#### Login

Вход в систему. Начало сеанса: генерация SID, хранение его в куках.

#### Logout

Выход из системы. Окончание сеанса: чистка SID и кук.

#### ConfirmUser

#### RejectUser

## GridController

#### GetSubTable

#### CreateGrid

#### EditGridRow

#### JsonToExcel

#### GetEventTypesByEventId

Возвращает названия типов мероприятий по ID мероприятия.

#### ImportForms

Копирование форм из последнесозданных мероприятий указанных типов мероприятий.

#### GetPersonsByEventId

Список значений об участниках мероприятия по параметрам анкет.

#### GetParamsByEventId

Список параметров по идентификатору мероприятия.

#### Load

Подгрузка данных в грид для любой таблицы БД.

## IndexController

#### Index

Загрузка стартовой страницы сайта. Анонс мероприятий.

#### Init

Сброс БД. Загрузка тестовых данных.

#### LoadContestsFromCats

Загрузка в БД незавершенных мероприятий из CATS.

#### CreateRegistrationEvent

Создание мероприятия, соответсвующих форм и полей для регистрации пользователей в системе.

## BlankController

#### GetPersonRequestFromGroup

#### GetPersonRequest

#### EditParams

#### GetEditHistoryData

#### GetHistoryRequest

#### GetListHistoryEvents

#### GetRequest

## Handler

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

## GroupController

#### Register

#### ConfirmInvitationToGroup

#### RejectInvitationToGroup

#### IsRegGroup

#### AddPerson

## UserController

#### ShowCabinet

#### Login

#### CheckEnable

#### CheckSession

#### ResetPassword

#### ConfirmOrRejectPersonRequest

#### SendEmailWelcomeToPofile

[scheme-1]: ../docs/img/controllers.png "controllers"
[scheme-2]: ../docs/img/controllers-update-1.png "controllers"
