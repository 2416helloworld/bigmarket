SET NAMES utf8mb4;

CREATE DATABASE IF NOT EXISTS big_market DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE big_market;

-- ===================== 奖品 =====================
CREATE TABLE IF NOT EXISTS award (
  id BIGINT NOT NULL AUTO_INCREMENT,
  award_id INT NOT NULL COMMENT '奖品ID',
  award_key VARCHAR(64) NOT NULL COMMENT '奖品发放类型',
  award_config VARCHAR(255) DEFAULT NULL COMMENT '奖品配置',
  award_desc VARCHAR(255) DEFAULT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_award_id (award_id)
) ENGINE=InnoDB COMMENT='奖品配置';

-- ===================== 日常行为返利 =====================
CREATE TABLE IF NOT EXISTS daily_behavior_rebate (
  id BIGINT NOT NULL AUTO_INCREMENT,
  behavior_type VARCHAR(16) NOT NULL COMMENT 'sign 签到',
  rebate_desc VARCHAR(128) NOT NULL,
  rebate_type VARCHAR(16) NOT NULL COMMENT 'sku / integral',
  rebate_config VARCHAR(64) NOT NULL,
  state VARCHAR(8) NOT NULL DEFAULT 'open',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_behavior_rebate (behavior_type, rebate_type)
) ENGINE=InnoDB COMMENT='日常行为返利';

-- ===================== 活动 =====================
CREATE TABLE IF NOT EXISTS raffle_activity (
  id BIGINT NOT NULL AUTO_INCREMENT,
  activity_id BIGINT NOT NULL,
  activity_name VARCHAR(64) NOT NULL,
  activity_desc VARCHAR(128) DEFAULT NULL,
  begin_date_time DATETIME NOT NULL,
  end_date_time DATETIME NOT NULL,
  strategy_id BIGINT NOT NULL,
  state VARCHAR(8) NOT NULL DEFAULT 'create' COMMENT 'create/open/close',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_activity_id (activity_id),
  KEY idx_strategy_id (strategy_id)
) ENGINE=InnoDB COMMENT='抽奖活动';

CREATE TABLE IF NOT EXISTS raffle_activity_count (
  id BIGINT NOT NULL AUTO_INCREMENT,
  activity_count_id BIGINT NOT NULL,
  total_count INT NOT NULL,
  day_count INT NOT NULL,
  month_count INT NOT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_activity_count_id (activity_count_id)
) ENGINE=InnoDB COMMENT='活动次数配置';

CREATE TABLE IF NOT EXISTS raffle_activity_sku (
  id BIGINT NOT NULL AUTO_INCREMENT,
  sku BIGINT NOT NULL,
  activity_id BIGINT NOT NULL,
  activity_count_id BIGINT NOT NULL,
  stock_count INT NOT NULL COMMENT '总库存',
  stock_count_surplus INT NOT NULL COMMENT '剩余库存',
  product_amount DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '积分价格',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_sku (sku),
  KEY idx_activity_id (activity_id)
) ENGINE=InnoDB COMMENT='活动商品SKU';

-- ===================== 策略 =====================
CREATE TABLE IF NOT EXISTS strategy (
  id BIGINT NOT NULL AUTO_INCREMENT,
  strategy_id BIGINT NOT NULL,
  strategy_desc VARCHAR(128) NOT NULL,
  rule_models VARCHAR(256) DEFAULT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_strategy_id (strategy_id)
) ENGINE=InnoDB COMMENT='抽奖策略';

CREATE TABLE IF NOT EXISTS strategy_award (
  id BIGINT NOT NULL AUTO_INCREMENT,
  strategy_id BIGINT NOT NULL,
  award_id INT NOT NULL,
  award_title VARCHAR(128) NOT NULL,
  award_subtitle VARCHAR(128) DEFAULT NULL,
  award_count INT NOT NULL DEFAULT 0,
  award_count_surplus INT NOT NULL DEFAULT 0,
  award_rate DECIMAL(8,6) NOT NULL DEFAULT 0,
  rule_models VARCHAR(256) DEFAULT NULL,
  sort INT NOT NULL DEFAULT 0,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_strategy_award (strategy_id, award_id)
) ENGINE=InnoDB COMMENT='策略奖品概率';

CREATE TABLE IF NOT EXISTS strategy_rule (
  id BIGINT NOT NULL AUTO_INCREMENT,
  strategy_id BIGINT NOT NULL,
  award_id INT DEFAULT NULL,
  rule_type TINYINT NOT NULL DEFAULT 1 COMMENT '1策略规则 2奖品规则',
  rule_model VARCHAR(32) NOT NULL,
  rule_value VARCHAR(256) NOT NULL,
  rule_desc VARCHAR(128) DEFAULT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_strategy_rule_model (strategy_id, rule_model)
) ENGINE=InnoDB COMMENT='策略规则';

-- ===================== 规则树 =====================
CREATE TABLE IF NOT EXISTS rule_tree (
  id BIGINT NOT NULL AUTO_INCREMENT,
  tree_id VARCHAR(32) NOT NULL,
  tree_name VARCHAR(64) NOT NULL,
  tree_desc VARCHAR(128) DEFAULT NULL,
  tree_node_rule_key VARCHAR(32) NOT NULL COMMENT '根节点',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_tree_id (tree_id)
) ENGINE=InnoDB COMMENT='规则树';

CREATE TABLE IF NOT EXISTS rule_tree_node (
  id BIGINT NOT NULL AUTO_INCREMENT,
  tree_id VARCHAR(32) NOT NULL,
  rule_key VARCHAR(32) NOT NULL,
  rule_desc VARCHAR(64) DEFAULT NULL,
  rule_value VARCHAR(128) DEFAULT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_tree_id (tree_id)
) ENGINE=InnoDB COMMENT='规则树节点';

CREATE TABLE IF NOT EXISTS rule_tree_node_line (
  id BIGINT NOT NULL AUTO_INCREMENT,
  tree_id VARCHAR(32) NOT NULL,
  rule_node_from VARCHAR(32) NOT NULL,
  rule_node_to VARCHAR(32) NOT NULL,
  rule_limit_type VARCHAR(16) NOT NULL DEFAULT 'EQUAL',
  rule_limit_value VARCHAR(16) NOT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_tree_id (tree_id)
) ENGINE=InnoDB COMMENT='规则树连线';

-- ===================== 用户活动账户 =====================
CREATE TABLE IF NOT EXISTS raffle_activity_account (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  activity_id BIGINT NOT NULL,
  total_count INT NOT NULL DEFAULT 0,
  total_count_surplus INT NOT NULL DEFAULT 0,
  day_count INT NOT NULL DEFAULT 0,
  day_count_surplus INT NOT NULL DEFAULT 0,
  month_count INT NOT NULL DEFAULT 0,
  month_count_surplus INT NOT NULL DEFAULT 0,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_activity (user_id, activity_id)
) ENGINE=InnoDB COMMENT='用户活动总账户';

CREATE TABLE IF NOT EXISTS raffle_activity_account_day (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  activity_id BIGINT NOT NULL,
  day DATE NOT NULL,
  day_count INT NOT NULL DEFAULT 0,
  day_count_surplus INT NOT NULL DEFAULT 0,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_activity_day (user_id, activity_id, day)
) ENGINE=InnoDB COMMENT='用户活动日账户';

CREATE TABLE IF NOT EXISTS raffle_activity_account_month (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  activity_id BIGINT NOT NULL,
  month VARCHAR(8) NOT NULL,
  month_count INT NOT NULL DEFAULT 0,
  month_count_surplus INT NOT NULL DEFAULT 0,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_activity_month (user_id, activity_id, month)
) ENGINE=InnoDB COMMENT='用户活动月账户';

-- ===================== 订单 =====================
CREATE TABLE IF NOT EXISTS raffle_activity_order (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  sku BIGINT NOT NULL,
  activity_id BIGINT NOT NULL,
  activity_name VARCHAR(64) NOT NULL,
  strategy_id BIGINT NOT NULL,
  order_id VARCHAR(32) NOT NULL,
  order_time DATETIME NOT NULL,
  total_count INT NOT NULL,
  day_count INT NOT NULL,
  month_count INT NOT NULL,
  pay_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  state VARCHAR(16) NOT NULL COMMENT 'completed/wait_pay/close',
  out_business_no VARCHAR(64) NOT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_out_business_no (user_id, out_business_no),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='活动充值/兑换订单';

CREATE TABLE IF NOT EXISTS user_raffle_order (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  activity_id BIGINT NOT NULL,
  strategy_id BIGINT NOT NULL,
  order_id VARCHAR(32) NOT NULL,
  order_time DATETIME NOT NULL,
  order_state VARCHAR(16) NOT NULL DEFAULT 'create',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_order_id (user_id, order_id),
  KEY idx_user_strategy_time (user_id, strategy_id, order_time)
) ENGINE=InnoDB COMMENT='用户抽奖单';

CREATE TABLE IF NOT EXISTS user_award_record (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  activity_id BIGINT NOT NULL,
  strategy_id BIGINT NOT NULL,
  order_id VARCHAR(32) NOT NULL,
  award_id INT NOT NULL,
  award_title VARCHAR(128) NOT NULL,
  award_config VARCHAR(255) DEFAULT NULL,
  award_state VARCHAR(16) NOT NULL DEFAULT 'create',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_order (user_id, order_id),
  KEY idx_user_activity (user_id, activity_id)
) ENGINE=InnoDB COMMENT='用户中奖记录';

CREATE TABLE IF NOT EXISTS user_behavior_rebate_order (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  order_id VARCHAR(32) NOT NULL,
  behavior_type VARCHAR(16) NOT NULL,
  rebate_desc VARCHAR(128) NOT NULL,
  rebate_type VARCHAR(16) NOT NULL,
  rebate_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  biz_id VARCHAR(64) NOT NULL,
  out_business_no VARCHAR(64) NOT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_out_business (user_id, out_business_no),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='用户行为返利订单';

-- ===================== 积分 =====================
CREATE TABLE IF NOT EXISTS user_credit_account (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  available_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  account_status VARCHAR(16) NOT NULL DEFAULT 'open',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_id (user_id)
) ENGINE=InnoDB COMMENT='用户积分账户';

CREATE TABLE IF NOT EXISTS user_credit_order (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  order_id VARCHAR(32) NOT NULL,
  trade_name VARCHAR(32) NOT NULL,
  trade_type VARCHAR(16) NOT NULL COMMENT 'forward/reverse',
  trade_amount DECIMAL(12,2) NOT NULL,
  out_business_no VARCHAR(64) NOT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_user_out_business (user_id, out_business_no),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='用户积分交易单';

-- ===================== 异步任务 =====================
CREATE TABLE IF NOT EXISTS task (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(32) NOT NULL,
  topic VARCHAR(64) NOT NULL,
  message_id VARCHAR(64) NOT NULL,
  message TEXT NOT NULL,
  state VARCHAR(16) NOT NULL DEFAULT 'create' COMMENT 'create/completed/fail',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_state_create_time (state, create_time),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='MQ任务补偿表';

-- ===================== 管理端账号 =====================
CREATE TABLE IF NOT EXISTS admin_user (
  id BIGINT NOT NULL AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL,
  password VARCHAR(128) NOT NULL COMMENT 'sha256',
  nickname VARCHAR(64) DEFAULT NULL,
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_username (username)
) ENGINE=InnoDB COMMENT='运营后台账号';
