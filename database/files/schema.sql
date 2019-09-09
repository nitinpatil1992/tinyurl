CREATE DATABASE tinyurl;
USE tinyurl;

DROP TABLE IF EXISTS `tiny_urls`;

CREATE TABLE `tiny_urls` (
    `short_url` char(13) NOT NULL UNIQUE,
    `long_url` varchar(255) NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`short_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `url_visits`;

CREATE TABLE `tiny_url_visits` (
    `short_url` char(13) NOT NULL,
    `visited_at` timestamp NOT NULL DEFAULT current_timestamp,
    INDEX (`short_url`),
    FOREIGN KEY (short_url) REFERENCES tiny_urls(short_url)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
