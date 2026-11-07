USE big_market;

INSERT INTO award (award_id, award_key, award_config, award_desc) VALUES
(101,'user_credit_random','1,100','用户随机积分'),
(102,'openai_use_count','5','OpenAI 增加使用次数'),
(103,'openai_use_count','10','OpenAI 增加使用次数'),
(104,'openai_use_count','20','OpenAI 增加使用次数'),
(105,'openai_model','gpt-4','OpenAI 增加模型'),
(106,'openai_model','dall-e-2','OpenAI 增加模型'),
(107,'openai_model','dall-e-3','OpenAI 增加模型'),
(108,'openai_use_count','100','OpenAI 增加使用次数'),
(109,'openai_model','gpt-4,dall-e-2,dall-e-3','OpenAI 增加模型'),
(100,'user_credit_blacklist','1','黑名单积分')
ON DUPLICATE KEY UPDATE award_key=VALUES(award_key), award_config=VALUES(award_config), award_desc=VALUES(award_desc);

INSERT INTO daily_behavior_rebate (behavior_type, rebate_desc, rebate_type, rebate_config, state) VALUES
('sign','签到返利-sku额度','sku','9014','open'),
('sign','签到返利-积分','integral','10','open')
ON DUPLICATE KEY UPDATE rebate_type=VALUES(rebate_type), rebate_config=VALUES(rebate_config), state=VALUES(state);

INSERT INTO strategy (strategy_id, strategy_desc, rule_models) VALUES
(100001,'抽奖策略','rule_blacklist,rule_weight'),
(100002,'抽奖策略-非完整1概率',NULL),
(100003,'抽奖策略-验证lock',NULL),
(100004,'抽奖策略-随机抽奖',NULL),
(100005,'抽奖策略-测试概率计算',NULL),
(100006,'抽奖策略-规则树','rule_blacklist,rule_weight')
ON DUPLICATE KEY UPDATE strategy_desc=VALUES(strategy_desc), rule_models=VALUES(rule_models);

INSERT INTO strategy_award (strategy_id, award_id, award_title, award_subtitle, award_count, award_count_surplus, award_rate, rule_models, sort) VALUES
(100001,101,'随机积分',NULL,80000,79998,0.300000,'tree_luck_award',1),
(100001,102,'5次使用',NULL,10000,9999,0.200000,'tree_luck_award',2),
(100001,103,'10次使用',NULL,5000,4998,0.200000,'tree_luck_award',3),
(100001,104,'20次使用',NULL,4000,3999,0.100000,'tree_luck_award',4),
(100001,105,'增加gpt-4对话模型',NULL,600,600,0.100000,'tree_luck_award',5),
(100001,106,'增加dall-e-2画图模型',NULL,200,200,0.050000,'tree_luck_award',6),
(100001,107,'增加dall-e-3画图模型','抽奖1次后解锁',200,200,0.040000,'tree_luck_award',7),
(100001,108,'增加100次使用','抽奖2次后解锁',199,199,0.009900,'tree_luck_award',8),
(100001,109,'解锁全部模型','抽奖6次后解锁',1,1,0.000100,'tree_luck_award',9),
(100002,101,'随机积分',NULL,1,1,0.500000,'tree_luck_award',1),
(100002,102,'5次使用',NULL,1,1,0.100000,'tree_luck_award',2),
(100002,106,'增加dall-e-2画图模型',NULL,1,1,0.010000,'tree_luck_award',3),
(100003,107,'增加dall-e-3画图模型','抽奖1次后解锁',200,200,0.040000,'tree_luck_award',7),
(100003,108,'增加100次使用','抽奖2次后解锁',199,199,0.009900,'tree_luck_award',8),
(100003,109,'解锁全部模型','抽奖6次后解锁',1,1,0.000100,'tree_luck_award',9),
(100004,109,'解锁全部模型','抽奖6次后解锁',1,1,1.000000,'tree_luck_award',9),
(100005,101,'随机积分',NULL,80000,80000,0.030000,'tree_luck_award',1),
(100005,102,'随机积分',NULL,80000,80000,0.030000,'tree_luck_award',1),
(100005,103,'随机积分',NULL,80000,80000,0.030000,'tree_luck_award',1),
(100005,104,'随机积分',NULL,80000,80000,0.030000,'tree_luck_award',1),
(100005,105,'随机积分',NULL,80000,80000,0.001000,'tree_luck_award',1),
(100006,101,'随机积分',NULL,10000,9982,0.790000,'tree_luck_award',1),
(100006,102,'OpenAI会员卡',NULL,100,15,0.030000,'tree_luck_award',2),
(100006,103,'支付优惠券',NULL,100,38,0.030000,'tree_luck_award',3),
(100006,104,'小米台灯',NULL,100,28,0.030000,'tree_luck_award',4),
(100006,105,'小米su7周体验','抽奖3次后解锁',100,37,0.030000,'tree_lock_3',5),
(100006,106,'轻奢办公椅','抽奖2次后解锁',100,22,0.030000,'tree_lock_2',6),
(100006,107,'小霸王游戏机','抽奖1次后解锁',100,21,0.030000,'tree_lock_1',7),
(100006,108,'暴走玩偶',NULL,100,25,0.030000,'tree_luck_award',8);

INSERT INTO strategy_rule (strategy_id, award_id, rule_type, rule_model, rule_value, rule_desc) VALUES
(100001,NULL,1,'rule_weight','60:102,103,104,105 200:106,107 1000:105','权重抽奖范围'),
(100001,NULL,1,'rule_blacklist','101:user001,user002,user003','黑名单抽奖，积分兜底'),
(100006,NULL,1,'rule_weight','10:102,103 70:106,107 1000:104,105','权重抽奖范围'),
(100006,NULL,1,'rule_blacklist','101:user001,user002,user003','黑名单抽奖，积分兜底')
ON DUPLICATE KEY UPDATE rule_value=VALUES(rule_value), rule_desc=VALUES(rule_desc);

INSERT INTO rule_tree (tree_id, tree_name, tree_desc, tree_node_rule_key) VALUES
('tree_lock_1','规则树','抽奖1次后解锁','rule_lock'),
('tree_lock_2','规则树','抽奖2次后解锁','rule_lock'),
('tree_lock_3','规则树','抽奖3次后解锁','rule_lock'),
('tree_luck_award','规则树-兜底奖励','兜底奖励','rule_stock')
ON DUPLICATE KEY UPDATE tree_name=VALUES(tree_name), tree_desc=VALUES(tree_desc);

INSERT INTO rule_tree_node (tree_id, rule_key, rule_desc, rule_value) VALUES
('tree_lock_1','rule_lock','限定用户已完成N次抽奖后解锁','1'),
('tree_lock_1','rule_luck_award','兜底奖品随机积分','101:1,100'),
('tree_lock_1','rule_stock','库存扣减规则',NULL),
('tree_lock_2','rule_lock','限定用户已完成N次抽奖后解锁','2'),
('tree_lock_2','rule_luck_award','兜底奖品随机积分','101:1,100'),
('tree_lock_2','rule_stock','库存扣减规则',NULL),
('tree_lock_3','rule_lock','限定用户已完成N次抽奖后解锁','3'),
('tree_lock_3','rule_luck_award','兜底奖品随机积分','101:1,100'),
('tree_lock_3','rule_stock','库存扣减规则',NULL),
('tree_luck_award','rule_stock','库存扣减规则',NULL),
('tree_luck_award','rule_luck_award','兜底奖品随机积分','101:1,100');

INSERT INTO rule_tree_node_line (tree_id, rule_node_from, rule_node_to, rule_limit_type, rule_limit_value) VALUES
('tree_lock_1','rule_lock','rule_stock','EQUAL','ALLOW'),
('tree_lock_1','rule_lock','rule_luck_award','EQUAL','TAKE_OVER'),
('tree_lock_1','rule_stock','rule_luck_award','EQUAL','ALLOW'),
('tree_lock_2','rule_lock','rule_stock','EQUAL','ALLOW'),
('tree_lock_2','rule_lock','rule_luck_award','EQUAL','TAKE_OVER'),
('tree_lock_2','rule_stock','rule_luck_award','EQUAL','ALLOW'),
('tree_lock_3','rule_lock','rule_stock','EQUAL','ALLOW'),
('tree_lock_3','rule_lock','rule_luck_award','EQUAL','TAKE_OVER'),
('tree_lock_3','rule_stock','rule_luck_award','EQUAL','ALLOW'),
('tree_luck_award','rule_stock','rule_luck_award','EQUAL','ALLOW');

INSERT INTO raffle_activity_count (activity_count_id, total_count, day_count, month_count) VALUES
(11101,150,150,150),
(11102,50,50,50),
(11103,5,5,5),
(11104,1,1,1),
(11105,10,10,10)
ON DUPLICATE KEY UPDATE total_count=VALUES(total_count), day_count=VALUES(day_count), month_count=VALUES(month_count);

INSERT INTO raffle_activity (activity_id, activity_name, activity_desc, begin_date_time, end_date_time, strategy_id, state) VALUES
(100301,'测试活动','小米 SU7 抽奖活动','2024-03-09 10:15:10','2034-03-09 10:15:10',100006,'open'),
(100302,'积分抽奖活动','签到得积分，积分兑换次数','2024-06-01 00:00:00','2034-06-01 00:00:00',100001,'open')
ON DUPLICATE KEY UPDATE activity_name=VALUES(activity_name), activity_desc=VALUES(activity_desc), strategy_id=VALUES(strategy_id), state=VALUES(state);

INSERT INTO raffle_activity_sku (sku, activity_id, activity_count_id, stock_count, stock_count_surplus, product_amount) VALUES
(9011,100301,11101,100000,99999,200.00),
(9012,100301,11102,100000,99999,100.00),
(9013,100301,11103,100000,99999,20.00),
(9014,100301,11104,100000,99999,5.00),
(9021,100302,11105,10000,10000,10.00)
ON DUPLICATE KEY UPDATE product_amount=VALUES(product_amount);

INSERT INTO admin_user (username, password, nickname) VALUES
('admin','240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9','运营管理员')
ON DUPLICATE KEY UPDATE nickname=VALUES(nickname);

INSERT INTO user_credit_account (user_id, total_amount, available_amount, account_status) VALUES
('user001',1000,1000,'open'),
('demo',500,500,'open')
ON DUPLICATE KEY UPDATE available_amount=VALUES(available_amount);

-- 演示账号 demo 预置 10 次活动额度，便于首次访问直接体验抽奖
INSERT INTO raffle_activity_account (user_id, activity_id, total_count, total_count_surplus, day_count, day_count_surplus, month_count, month_count_surplus) VALUES
('demo',100301,10,10,10,10,10,10)
ON DUPLICATE KEY UPDATE total_count_surplus=VALUES(total_count_surplus);
INSERT INTO raffle_activity_account_day (user_id, activity_id, day, day_count, day_count_surplus) VALUES
('demo',100301,DATE_FORMAT(NOW(),'%Y-%m-%d'),10,10)
ON DUPLICATE KEY UPDATE day_count_surplus=VALUES(day_count_surplus);
INSERT INTO raffle_activity_account_month (user_id, activity_id, month, month_count, month_count_surplus) VALUES
('demo',100301,DATE_FORMAT(NOW(),'%Y-%m'),10,10)
ON DUPLICATE KEY UPDATE month_count_surplus=VALUES(month_count_surplus);
