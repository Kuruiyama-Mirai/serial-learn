CREATE TABLE `save_data` (
  `id` int NOT NULL AUTO_INCREMENT,
  `deviceId` bigint NOT NULL COMMENT '设备ID',
  `temperature` float DEFAULT NULL COMMENT '温度',
  `temperatureAlarmStatus` int DEFAULT NULL COMMENT '温度报警状态',
  `diStatus` int DEFAULT NULL COMMENT 'DI状态',
  `soepointer` int DEFAULT NULL COMMENT 'SOE事件数',
  `voltageImbalance` float DEFAULT NULL COMMENT '电压不平衡度',
  `currentImbalance` float DEFAULT NULL COMMENT '电流不平衡度',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `index_deviceId` (`deviceId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;