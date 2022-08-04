CREATE DATABASE IF NOT EXISTS demo DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

use demo;

create table user
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime        null,
    updated_at  datetime        null,
    deleted_at  datetime        null,
    name        char(50)        not null comment '用户名',
    password    char(100)       not null comment '密码',
    email       char(50)        not null comment '邮件',
    phone       bigint unsigned not null comment '手机号码',
    age         tinyint         not null comment '年龄',
    gender      tinyint         not null comment '性别，1:男，2:女，3:未知',
    status      tinyint         not null comment '账号状态，1:未激活，2:已激活，3:封禁',
    login_state tinyint         not null comment '登录状态，1:未登录，2:已登录',
    constraint user_email_uindex
        unique (email)
);
