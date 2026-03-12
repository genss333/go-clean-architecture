/*
 Navicat Premium Dump SQL

 Source Server         : local-1
 Source Server Type    : PostgreSQL
 Source Server Version : 160013 (160013)
 Source Host           : localhost:5432
 Source Catalog        : cleanarch_db
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160013 (160013)
 File Encoding         : 65001

 Date: 12/03/2026 16:25:00
*/


-- ----------------------------
-- Table structure for departments
-- ----------------------------
DROP TABLE IF EXISTS "public"."departments";
CREATE TABLE "public"."departments" (
  "department_id" int4 NOT NULL DEFAULT nextval('departments_department_id_seq'::regclass),
  "department_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."departments" OWNER TO "root";
COMMENT ON COLUMN "public"."departments"."department_name" IS 'deaprtment name is unique';

-- ----------------------------
-- Primary Key structure for table departments
-- ----------------------------
ALTER TABLE "public"."departments" ADD CONSTRAINT "departments_pkey" PRIMARY KEY ("department_id", "department_name");
