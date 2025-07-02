/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE IF NOT EXISTS event_management;
USE event_management;

DROP TABLE IF EXISTS `attendee_status`;
CREATE TABLE `attendee_status` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `event_attendees`;
CREATE TABLE `event_attendees` (
  `event_id` int NOT NULL,
  `user_id` int NOT NULL,
  `status_id` int NOT NULL DEFAULT '1',
  UNIQUE KEY `event_user_unique` (`event_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `events`;
CREATE TABLE `events` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` text,
  `location` varchar(250) DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `created_by` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_public` tinyint(1) NOT NULL DEFAULT '0',
  `attendee_limit` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_events_created_by` (`created_by`),
  CONSTRAINT `fk_events_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `permission` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `description` varchar(100) NOT NULL DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `permission` (`permission`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
  `role_id` int NOT NULL,
  `permission_id` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_role_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`),
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb3;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  `first_name` varchar(50) NOT NULL,
  `last_name` varchar(50) NOT NULL,
  `role_id` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_users_role_id` (`role_id`),
  CONSTRAINT `fk_users_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3;

INSERT INTO `attendee_status` (`id`, `title`) VALUES
(1, 'Invited'),
(2, 'Accepted'),
(3, 'Rejected');

INSERT INTO `event_attendees` (`event_id`, `user_id`, `status_id`) VALUES
(2, 1, 0),
(9, 2, 0),
(10, 2, 3),
(11, 2, 0),
(12, 2, 0),
(13, 2, 0),
(14, 2, 0),
(15, 2, 0),
(16, 2, 0),
(18, 1, 0),
(19, 1, 1),
(20, 1, 1),
(21, 1, 1),
(22, 2, 1),
(23, 2, 1),
(25, 2, 1),
(26, 2, 1),
(27, 2, 1),
(32, 2, 2),
(33, 1, 2),
(33, 4, 2);

INSERT INTO `events` (`id`, `title`, `description`, `location`, `start_time`, `end_time`, `created_by`, `created_at`, `updated_at`, `is_public`, `attendee_limit`) VALUES
(1, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:11:38', '2025-05-29 02:11:38', 1, NULL),
(2, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:12:02', '2025-05-29 02:12:02', 1, 100),
(3, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:12:36', '2025-05-29 02:12:36', 1, 100),
(4, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:12:37', '2025-05-29 02:12:37', 1, 100),
(5, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:12:38', '2025-05-29 02:12:38', 1, 100),
(6, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:17:00', '2025-05-29 02:17:00', 0, 100),
(7, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 02:17:44', '2025-05-29 02:17:44', 0, 100),
(9, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 03:04:16', '2025-05-29 03:04:16', 0, 100),
(10, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 06:41:47', '2025-05-29 06:41:47', 0, 100),
(11, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 06:41:50', '2025-05-29 06:41:50', 0, 100),
(12, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 06:41:50', '2025-05-29 06:41:50', 0, 100),
(13, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 06:41:51', '2025-05-29 06:41:51', 0, 100),
(14, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-29 06:41:52', '2025-05-29 06:41:52', 0, 100),
(15, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 3, '2025-05-29 09:57:31', '2025-05-29 09:57:31', 0, 100),
(16, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 3, '2025-05-29 09:57:41', '2025-05-29 09:57:41', 0, 100),
(17, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 3, '2025-05-29 09:57:52', '2025-05-29 09:57:52', 0, 100),
(18, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 3, '2025-05-29 09:59:45', '2025-05-29 09:59:45', 0, 100),
(19, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 17:17:06', '2025-05-30 17:17:06', 0, NULL),
(20, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 17:18:45', '2025-05-30 17:18:45', 0, NULL),
(21, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 17:20:56', '2025-05-30 17:20:56', 0, NULL),
(22, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 17:48:37', '2025-05-30 17:48:37', 0, NULL),
(23, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-06-14 11:00:00', 1, '2025-05-30 17:48:51', '2025-05-31 04:56:28', 1, NULL),
(24, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 17:50:59', '2025-05-30 17:50:59', 1, NULL),
(25, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-30 18:12:37', '2025-05-30 18:12:37', 0, NULL),
(26, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-31 04:11:42', '2025-05-31 04:11:42', 0, NULL),
(27, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-05-14 10:00:00', '2025-05-14 11:00:00', 1, '2025-05-31 04:17:29', '2025-05-31 04:17:29', 0, NULL),
(28, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-06-10 10:00:00', '2025-06-10 11:00:00', 1, '2025-05-31 09:26:50', '2025-05-31 09:26:50', 0, 0),
(29, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-06-10 10:00:00', '2025-06-10 11:00:00', 1, '2025-05-31 09:27:08', '2025-05-31 09:27:08', 0, 100),
(30, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-06-10 10:00:00', '2025-06-10 11:00:00', 1, '2025-05-31 09:27:19', '2025-05-31 09:27:19', 1, 100),
(32, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-06-10 10:00:00', '2025-06-10 11:00:00', 1, '2025-05-31 09:44:22', '2025-05-31 09:44:22', 0, 100),
(33, 'Project Structure kickoff', 'Introduction to the course. Incebreaking seesion for trainer and students', 'Office', '2025-06-10 10:00:00', '2025-06-10 11:00:00', 1, '2025-05-31 11:00:31', '2025-05-31 11:08:37', 1, 1);

INSERT INTO `permissions` (`id`, `permission`, `description`, `created_at`, `updated_at`) VALUES
(1, 'user.create', 'Permission to create a new user', '2025-05-28 18:02:52', NULL),
(2, 'user.update', 'Permission to update an existing user', '2025-05-28 18:02:52', NULL),
(3, 'user.fetch', 'Permission to fetch a specific user', '2025-05-28 18:02:52', NULL),
(4, 'user.list', 'Permission to list users', '2025-05-28 18:02:52', NULL),
(5, 'user.delete', 'Permission to delete a user', '2025-05-28 18:02:52', NULL),
(6, 'event.create', 'Permission to create a new event', '2025-05-28 18:02:52', NULL),
(7, 'event.update', 'Permission to update an existing event', '2025-05-28 18:02:52', NULL),
(8, 'event.fetch', 'Permission to fetch a specific event', '2025-05-28 18:02:52', NULL),
(9, 'event.list', 'Permission to list events', '2025-05-28 18:02:52', NULL),
(10, 'event.delete', 'Permission to delete an event', '2025-05-28 18:02:52', NULL),
(11, 'user.fetchAllUserAsAttendee', 'permission to fetch all the user as attendee', '2025-05-29 11:12:27', '2025-05-29 11:18:57'),
(12, 'user.listAttendee', 'permission to list attendee', '2025-05-29 11:18:19', NULL),
(13, 'event.fetchAllEvent', 'admin permission for fetch event', '2025-05-29 12:33:35', NULL),
(14, 'event.fetchOwnEvent', 'fetch event created by own', '2025-05-29 12:34:18', NULL),
(15, 'event.fetchInvitedEvent', 'fetch invited event', '2025-05-29 12:35:19', NULL);

INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES
(1, 1, '2025-05-28 18:02:52', NULL),
(1, 2, '2025-05-28 18:02:52', NULL),
(1, 3, '2025-05-28 18:02:52', NULL),
(1, 4, '2025-05-28 18:02:52', NULL),
(1, 5, '2025-05-28 18:02:52', NULL),
(1, 6, '2025-05-28 18:02:52', NULL),
(1, 7, '2025-05-28 18:02:52', NULL),
(1, 8, '2025-05-28 18:02:52', NULL),
(1, 9, '2025-05-28 18:02:52', NULL),
(1, 10, '2025-05-28 18:02:52', NULL),
(1, 11, '2025-05-29 11:12:56', NULL),
(1, 13, '2025-05-29 12:40:48', NULL),
(2, 3, '2025-05-28 18:02:52', NULL),
(2, 4, '2025-05-28 18:02:52', NULL),
(2, 6, '2025-05-28 18:02:52', NULL),
(2, 7, '2025-05-28 18:02:52', NULL),
(2, 8, '2025-05-28 18:02:52', NULL),
(2, 9, '2025-05-28 18:02:52', NULL),
(2, 10, '2025-05-28 18:02:52', NULL),
(2, 14, '2025-05-29 12:40:48', NULL),
(3, 8, '2025-05-28 18:02:52', NULL),
(3, 9, '2025-05-28 18:02:52', NULL),
(3, 15, '2025-05-29 12:40:48', NULL);

INSERT INTO `roles` (`id`, `name`) VALUES
(1, 'ADMIN'),
(2, 'MANAGER'),
(3, 'ATTENDEE');

INSERT INTO `users` (`id`, `email`, `password`, `first_name`, `last_name`, `role_id`, `created_at`, `updated_at`) VALUES
(1, 'admin@vivasoftltd.com', '$2a$10$2fBRiXac/mWv9m1n891zv.K1ooO1ItZtArxGqpO5qFEX6xgtgrDzu', 'Abdul', 'Mukit', 1, '2025-05-28 18:02:52', NULL),
(2, 'user1@example.com', '$2a$10$ktRzNbzas/SmQpXJEH6MEOWkfOmPszxOPWo6wuD2eQ1y3u.hEQ.pu', 'User1', 'Last', 3, '2025-05-29 01:15:13', '2025-05-31 06:26:18'),
(3, 'manager@vivasoftltd.com', '$2a$10$2fBRiXac/mWv9m1n891zv.K1ooO1ItZtArxGqpO5qFEX6xgtgrDzu', 'Mostakim', 'Billah', 2, '2025-05-28 18:02:52', NULL),
(4, 'user2@example.com', '$2a$10$ktRzNbzas/SmQpXJEH6MEOWkfOmPszxOPWo6wuD2eQ1y3u.hEQ.pu', 'User2', 'Last', 3, '2025-05-29 01:15:13', '2025-05-31 06:26:18'),
(5, 'user3@example.com', '$2a$10$ktRzNbzas/SmQpXJEH6MEOWkfOmPszxOPWo6wuD2eQ1y3u.hEQ.pu', 'User3', 'Last', 3, '2025-05-29 01:15:13', '2025-05-31 06:26:18');



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;