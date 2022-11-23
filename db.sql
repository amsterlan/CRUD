CREATE TABLE `employees` (
                         `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(30) NOT NULL,
                         `email` varchar(30) NOT NULL,
                         `salary` decimal(20,2) NOT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;