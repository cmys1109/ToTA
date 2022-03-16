# ToTA

这种形式的留言板我之前玩过，也碰到过一些有趣  
碰到过一些巧合，让我想起过从前的某人某事

在2022/03/12的晚饭过后，我闲着无聊就突然想起了这个，
就顺手写了一下，因为之前对MySQL基本上没有了解，
所以就把这个当作是自己学习的一个小项目，当然有很多不足与漏洞，
我并没有完整MySQL。并且在这个项目里，有些地方我也没有接触过，
属于第一次使用  

本项目的配置文件应该是一眼就懂的，代码没啥注释，但是功能都很基础，
看一下接口应该就能知道干啥用的  

后端 我觉得需要增加防止恶意访问的功能

前端 我根本都不了解，所以就就随便用HTML写了一下，十分简陋
仅仅满足最基本的使用需求  

## MySQL DDL:
```MySQL DDL
create table ToTA
(
    id   int auto_increment
        primary key,
    name char(20) not null,
    text text     not null,
    time datetime not null,
    ip   char(20) null,
    constraint ToTA_id_uindex
        unique (id)
);
```


```
最后祝:Can say anything to sb,那个TA
```