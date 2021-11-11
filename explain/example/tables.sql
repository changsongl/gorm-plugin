CREATE TABLE `test`.`explain_table` (
`id` INT NOT NULL AUTO_INCREMENT,
`room_id` INT UNSIGNED NOT NULL,
`room_name` VARCHAR(45) NOT NULL,
PRIMARY KEY (`id`),
INDEX `idx_room_id` (`room_id` ASC));

INSERT INTO `test`.`explain_table` (`id`, `room_id`, `room_name`) VALUES (NULL, '1', 'room_name_1');
INSERT INTO `test`.`explain_table` (`id`, `room_id`, `room_name`) VALUES (NULL, '2', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('3', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('4', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('5', 'room_name3');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('5', 'room_name5');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('6', 'room_name9');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('7', 'room_name1');
