CREATE TABLE `statistics`
(
    `id_slot`     int NOT NULL COMMENT 'id слота',
    `id_banner`   int NOT NULL COMMENT 'id баннера',
    `id_group`    int NOT NULL COMMENT 'id социальной группы',
    `count_click` int DEFAULT '0' COMMENT 'кол-во кликов',
    `count_show`  int DEFAULT '0' COMMENT 'кол-во просмотров',
    PRIMARY KEY (`id_slot`, `id_group`, `id_banner`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `rotation`
(
    `id_banner` int NOT NULL COMMENT 'id баннера',
    `id_slot`   int NOT NULL COMMENT 'id слота',
    `status`    tinyint DEFAULT NULL COMMENT 'статус, 0 - не удален, 1 - удален',
    PRIMARY KEY (`id_banner`, `id_slot`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
