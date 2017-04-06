--- Big rank by asin!
  CREATE TABLE `A1` (
  `id` varchar(255) NOT NULL,
  `day` varchar(255) NOT NULL,
  `bigname` varchar(255) DEFAULT NULL COMMENT '大类名',
  `title` TEXT COMMENT '商品标题',
  `rank` int(11) DEFAULT NULL COMMENT '大类排名',
  `price` float DEFAULT NULL,
  `sold` varchar(255) DEFAULT NULL COMMENT '自营',
  `ship` varchar(255) DEFAULT NULL COMMENT 'FBA',
  `score` float DEFAULT NULL COMMENT '打分',
  `reviews` int(11) DEFAULT NULL COMMENT '评论数',
  `createtime` varchar(255) DEFAULT NULL,
  `img` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`,`day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


  CREATE TABLE `A2` (
  `id` varchar(150) NOT NULL,
  `day` varchar(150) NOT NULL,
  `bigname` varchar(255) DEFAULT NULL COMMENT '大类名',
  `title` TEXT COMMENT '商品标题',
  `rank` int(11) DEFAULT NULL COMMENT '大类排名',
  `price` float DEFAULT NULL,
  `sold` varchar(255) DEFAULT NULL COMMENT '自营',
  `ship` varchar(255) DEFAULT NULL COMMENT 'FBA',
  `score` float DEFAULT NULL COMMENT '打分',
  `reviews` int(11) DEFAULT NULL COMMENT '评论数',
  `createtime` varchar(255) DEFAULT NULL,
  `img` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`,`day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;