--- smart_base
--- Category
CREATE TABLE `smart_category` (
  `id` varchar(100) NOT NULL,
  `url` varchar(255) DEFAULT NULL COMMENT '类目链接',
  `name` varchar(255) DEFAULT NULL COMMENT '类目名字',
  `level` tinyint(4) DEFAULT NULL COMMENT '类目级别',
  `pid` varchar(100) DEFAULT NULL COMMENT '父类id',
  `createtime` datetime DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime DEFAULT NULL COMMENT '更新时间',
  `isvalid` tinyint(4) DEFAULT '0' COMMENT '是否有效',
  `page` tinyint(4) DEFAULT '5' COMMENT '抓取页数',
  `database` varchar(255) DEFAULT NULL COMMENT '存储数据库',
  `col1` varchar(255) DEFAULT NULL COMMENT '预留字段',
  `col2` varchar(255) DEFAULT NULL,
  `col3` varchar(255) DEFAULT NULL,
  `bigpname` varchar(255) DEFAULT NULL COMMENT '大类名字',
  `bigpid` varchar(100) DEFAULT NULL COMMENT '大类ID',
  `ismall` tinyint(4) DEFAULT '0' COMMENT '是否最小类',
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_UNIQUE` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='类目';


--- ASIN Collect
CREATE TABLE `smart_asin` (
  `id` varchar(100) NOT NULL,
  `createtime` varchar(255) DEFAULT NULL COMMENT '添加时间',
  `updatetime` varchar(255) DEFAULT NULL COMMENT '更新时间',
  `category` varchar(255) DEFAULT NULL COMMENT "which category",
  `times` int(11) DEFAULT '0' COMMENT '重复次数',
  `isvalid` tinyint(4) DEFAULT '1' COMMENT "valid",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='Asin Big Data';


--- ASIN Collect
CREATE TABLE `smart_monitor` (
  `id` VARCHAR(30) NOT NULL,
  `redistype` VARCHAR(50) NOT NULL,
  `doing` INT NULL,
  `done` INT NULL,
  `updatetime` VARCHAR(45) NULL,
  `createtime` VARCHAR(45) NULL,
  PRIMARY KEY (`id`, `redistype`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;
