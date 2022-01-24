-- +goose Up
CREATE TABLE `countries` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `code` varchar(3)
);

CREATE TABLE `cases` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `country_id` int,
  `date_id` int,
  `confirmed` int,
  `deaths` int,
  `stringency_actual` decimal(5,2),
  `stringency` decimal(5,2)
);

CREATE TABLE `dates` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `date_value` date
);

ALTER TABLE `cases` ADD FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`);

ALTER TABLE `cases` ADD FOREIGN KEY (`date_id`) REFERENCES `dates` (`id`);

-- +goose Down
DROP TABLE countries;
DROP TABLE cases;
DROP TABLE dates;
