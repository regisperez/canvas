DROP DATABASE IF EXISTS `sketch_test`;
CREATE DATABASE IF NOT EXISTS `sketch_test`;
DROP DATABASE IF EXISTS `sketch`;
CREATE DATABASE IF NOT EXISTS `sketch`;
USE `sketch`;
DROP TABLE IF EXISTS `canvas`;
CREATE TABLE `canvas` (
                          `id` varchar(64) NOT NULL,
                          `drawing` varchar(1000) NOT NULL,
                          `creationdate` datetime NOT NULL,
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;