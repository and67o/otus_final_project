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



INSERT INTO rotation (id, id_banner, id_slot) VALUES (23, 1, 2);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (22, 1, 3);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (33, 2, 1);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (25, 2, 2);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (28, 2, 3);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (26, 3, 1);
INSERT INTO rotation (id, id_banner, id_slot) VALUES (27, 4, 3);