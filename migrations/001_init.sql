-- ★ 为什么第一行必须是 SET NAMES utf8mb4？
--   mysql 客户端的 default-character-set 默认是 auto，会按操作系统 locale 推断字符集；
--   docker exec 进容器时通常没有 LANG（POSIX/C locale），auto 会退化成 latin1。
--   此时文件里的 UTF-8 中文会被当成 latin1 解释，服务端再转成 utf8mb4 落库 → 双重编码乱码：
--     · Navicat（utf8mb4 连接）如实显示 → å•†å“   ← 看着是乱码
--     · mysql CLI（同样是 auto→latin1 连接）读出来时字节被"还原" → 看着正常（假象！）
--   把 SET NAMES 写进文件，正确的连接字符集就跟着文件走，不依赖调用方是否传 --default-character-set。
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `store_product_spu`
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`    BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `category_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属分类',
    `name`        VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '商品名称',
    `code`        VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '商品code',
    `img_url`     VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '商品图片',
    `description` VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '商品描述',
    `status`      TINYINT         NOT NULL DEFAULT 0 COMMENT '商品状态，0-下架，1-上架',
    `sort`        INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '商品排序，值越大，排序越靠前',
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`  DATETIME        NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_store_code` (`store_id`, `code`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='商品SPU';


CREATE TABLE IF NOT EXISTS `store_product_sku`
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`    BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `spu_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属spu_id',
    `name`        VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '商品名称',
    `code`        VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '商品code',
    `price`       INT             NOT NULL DEFAULT 0 COMMENT '售价，单位分',
    `cost_price`  INT             NOT NULL DEFAULT 0 COMMENT '成本价，单位分',
    `img_url`     VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '商品图片',
    `description` VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '商品描述',
    `status`      TINYINT         NOT NULL DEFAULT 0 COMMENT '商品状态，0-下架，1-上架',
    `sort`        INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '商品排序，值越大，排序越靠前',
    `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`  DATETIME        NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_store_code` (`store_id`, `code`),
    KEY `idx_store_spu` (`store_id`, `spu_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='商品SKU';


CREATE TABLE IF NOT EXISTS `store_product_inventory`
(
    `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`        BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `spu_id`          BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属spu_id（冗余字段）',
    `sku_id`          BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属sku_id',
    `total_stock`     BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际库存',
    `available_stock` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '可售库存',
    `locked_stock`    BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '锁定的库存',
    `warn_stock`      INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '预警库存',
    `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`      DATETIME        NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_store_sku` (`store_id`, `sku_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='商品库存';


CREATE TABLE IF NOT EXISTS `store_product_inventory_log`
(
    `id`           BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`     BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `spu_id`       BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '所属spu_id（冗余字段）',
    `sku_id`       BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '所属sku_id',
    `quantity`     BIGINT           NOT NULL DEFAULT 0 COMMENT '变动数量（正负）',
    `before_stock` BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '变动前的库存',
    `after_stock`  BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '变动后的库存',
    `biz_type`     TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '业务来源类型，1-下单/2-取消订单/3-退款/4-入库/5-出库',
    `biz_no`   VARCHAR(32)      NOT NULL DEFAULT '' COMMENT '业务单号',
    `created_at`   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`   DATETIME         NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    KEY `idx_store_sku` (`store_id`, `sku_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='库存流程';



CREATE TABLE IF NOT EXISTS `store_order`
(
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`        BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `order_no`        VARCHAR(32)      NOT NULL DEFAULT '' COMMENT '订单号',
    `user_id`         BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '下单的用户ID',
    `order_type`      TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单类型，0-普通/1-秒杀',
    `status`          TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单状态，0-未完成/1-已完成/2-取消中/3-已取消/4-退款中/5-已退款',
    `pay_status`      TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '支付状态，0-未支付/1-部分支付/2-支付完成',
    `pay_time`        DATETIME         NULL     DEFAULT NULL COMMENT '支付时间，NULL=未支付',
    `total_amount`    BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '订单总金额，单位分',
    `discount_amount` BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '折扣金额，单位分',
    `pay_amount`      BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '实际支付金额金额，单位分',
    `remark`          VARCHAR(128)     NOT NULL DEFAULT '' COMMENT '订单备注',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`      DATETIME         NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_store_order` (`store_id`, `order_no`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='订单主表';


CREATE TABLE IF NOT EXISTS `store_order_item`
(
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `store_id`        BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '店铺ID',
    `order_id`        BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '订单ID',
    `spu_id`          BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '所属spu_id',
    `sku_id`          BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '所属sku_id',
    `sku_code`        VARCHAR(32)      NOT NULL DEFAULT '' COMMENT '商品sku code',
    `sku_name`        VARCHAR(64)      NOT NULL DEFAULT '' COMMENT '商品sku名称-快照',
    `sku_img_url`     VARCHAR(128)     NOT NULL DEFAULT '' COMMENT '商品图片-快照',
    `price`           INT              NOT NULL DEFAULT 0 COMMENT '下单时的单价，单位分',
    `quantity`        INT              NOT NULL DEFAULT 0 COMMENT '数量',
    `total_amount`    BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '小计=price * quantity，单位分',
    `discount_amount` BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '折扣金额，单位分',
    `pay_amount`      BIGINT UNSIGNED  NOT NULL DEFAULT 0 COMMENT '实际支付金额金额，单位分',
    `status`          TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单状态，0-未完成/1-已完成/2-取消中/3-已取消/4-退款中/5-已退款',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`      DATETIME         NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
    PRIMARY KEY (`id`),
    KEY `idx_store_order` (`store_id`, `order_id`),
    KEY `idx_store_sku` (`store_id`, `sku_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='订单明细';


# # 暂定
# CREATE TABLE IF NOT EXISTS `store_order_cancel`
# (
#     `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
#     `store_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺ID',
#     `order_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单ID',
#     `cancel_time`   DATETIME        NULL     DEFAULT NULL COMMENT '取消时间',
#     `cancel_reason` VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '订单备注',
#     `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
#     `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
#     `deleted_at`    DATETIME        NULL     DEFAULT NULL COMMENT '删除时间，NULL=未删除',
#     PRIMARY KEY (`id`),
#     UNIQUE KEY `uk_store_order` (`store_id`, `order_id`)
# ) ENGINE = InnoDB
#   DEFAULT CHARSET = utf8mb4 COMMENT ='订单取消原因表';




