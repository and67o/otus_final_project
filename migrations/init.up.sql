CREATE TABLE `statistics`
(
    `id_slot`     int NOT NULL,
    `id_banner`   int NOT NULL,
    `id_group`    int NOT NULL,
    `count_click` int DEFAULT '0',
    `count_show`  int DEFAULT '0',
    PRIMARY KEY (`id_slot`, `id_group`, `id_banner`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `rotation`
(
    `id`        int NOT NULL AUTO_INCREMENT,
    `id_banner` int NOT NULL,
    `id_slot`   int NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `rotation_pk` (`id_banner`, `id_slot`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 30
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
