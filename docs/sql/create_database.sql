-- DROP SCHEMA store_go;

CREATE SCHEMA store_go AUTHORIZATION postgres;

-- DROP SEQUENCE application_contact_id_seq;

CREATE SEQUENCE application_contact_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE application_detail_id_seq;

CREATE SEQUENCE application_detail_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE application_id_seq;

CREATE SEQUENCE application_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE application_profile_history_id_seq;

CREATE SEQUENCE application_profile_history_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE application_version_id_seq;

CREATE SEQUENCE application_version_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE category_id_seq;

CREATE SEQUENCE category_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE category_type_id_seq;

CREATE SEQUENCE category_type_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE image_id_seq;

CREATE SEQUENCE image_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE integration_type_id_seq;

CREATE SEQUENCE integration_type_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE storage_object_id_seq;

CREATE SEQUENCE storage_object_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE terminal_model_id_seq;

CREATE SEQUENCE terminal_model_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE video_id_seq;

CREATE SEQUENCE video_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;-- store_go.application definition

-- Drop table

-- DROP TABLE application;

CREATE TABLE application (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	CONSTRAINT application_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_application_ativo ON store_go.application USING btree (ativo);
CREATE INDEX idx_application_created_at ON store_go.application USING btree (created_at);
CREATE INDEX idx_application_nome ON store_go.application USING btree (nome);


-- store_go.application_contact definition

-- Drop table

-- DROP TABLE application_contact;

CREATE TABLE application_contact (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	site text NOT NULL,
	email text NOT NULL,
	phone text NOT NULL,
	CONSTRAINT application_contact_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_application_contact_ativo ON store_go.application_contact USING btree (ativo);
CREATE INDEX idx_application_contact_created_at ON store_go.application_contact USING btree (created_at);
CREATE INDEX idx_application_contact_nome ON store_go.application_contact USING btree (nome);


-- store_go.application_detail definition

-- Drop table

-- DROP TABLE application_detail;

CREATE TABLE application_detail (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	descricao varchar(255) NULL,
	CONSTRAINT application_detail_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_application_detail_created_at ON store_go.application_detail USING btree (created_at);


-- store_go.category_type definition

-- Drop table

-- DROP TABLE category_type;

CREATE TABLE category_type (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	CONSTRAINT category_type_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_category_type_ativo ON store_go.category_type USING btree (ativo);
CREATE INDEX idx_category_type_created_at ON store_go.category_type USING btree (created_at);
CREATE INDEX idx_category_type_nome ON store_go.category_type USING btree (nome);


-- store_go.integration_type definition

-- Drop table

-- DROP TABLE integration_type;

CREATE TABLE integration_type (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	CONSTRAINT integration_type_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_integration_type_ativo ON store_go.integration_type USING btree (ativo);
CREATE INDEX idx_integration_type_created_at ON store_go.integration_type USING btree (created_at);
CREATE INDEX idx_integration_type_nome ON store_go.integration_type USING btree (nome);


-- store_go.storage_object definition

-- Drop table

-- DROP TABLE storage_object;

CREATE TABLE storage_object (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	"path" text NOT NULL,
	"name" varchar(40) NOT NULL,
	mime_type varchar(100) NOT NULL,
	status int2 NOT NULL,
	CONSTRAINT storage_object_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_storage_object_created_at ON store_go.storage_object USING btree (created_at);


-- store_go.terminal_model definition

-- Drop table

-- DROP TABLE terminal_model;

CREATE TABLE terminal_model (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	CONSTRAINT terminal_model_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_terminal_model_ativo ON store_go.terminal_model USING btree (ativo);
CREATE INDEX idx_terminal_model_created_at ON store_go.terminal_model USING btree (created_at);
CREATE INDEX idx_terminal_model_nome ON store_go.terminal_model USING btree (nome);


-- store_go.application_configuration definition

-- Drop table

-- DROP TABLE application_configuration;

CREATE TABLE application_configuration (
	application_id int8 NOT NULL,
	integration_type_id int8 NOT NULL,
	terminal_model_id int8 NOT NULL,
	CONSTRAINT fk_application_configuracoes_application FOREIGN KEY (application_id) REFERENCES application(id),
	CONSTRAINT fk_application_configuration_integration_type FOREIGN KEY (integration_type_id) REFERENCES integration_type(id),
	CONSTRAINT fk_application_configuration_terminal_model FOREIGN KEY (terminal_model_id) REFERENCES terminal_model(id)
);
CREATE UNIQUE INDEX uk_cfg_pk ON store_go.application_configuration USING btree (application_id, integration_type_id, terminal_model_id);
CREATE INDEX uk_cfg_rev ON store_go.application_configuration USING btree (terminal_model_id, integration_type_id, application_id);


-- store_go.application_profile_history definition

-- Drop table

-- DROP TABLE application_profile_history;

CREATE TABLE application_profile_history (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	application_contact_id int8 NULL,
	application_detail_id int8 NULL,
	review_at timestamptz NULL,
	production_at timestamptz NULL,
	deactivated_at timestamptz NULL,
	deactivation_cause varchar(255) NULL,
	CONSTRAINT application_profile_history_pkey PRIMARY KEY (id),
	CONSTRAINT fk_application_profile_history_application_contact FOREIGN KEY (application_contact_id) REFERENCES application_contact(id),
	CONSTRAINT fk_application_profile_history_application_detail FOREIGN KEY (application_detail_id) REFERENCES application_detail(id)
);
CREATE INDEX idx_application_profile_history_created_at ON store_go.application_profile_history USING btree (created_at);


-- store_go.application_profile_history_configuration definition

-- Drop table

-- DROP TABLE application_profile_history_configuration;

CREATE TABLE application_profile_history_configuration (
	application_profile_history_id int8 NOT NULL,
	application_id int8 NOT NULL,
	integration_type_id int8 NOT NULL,
	terminal_model_id int8 NOT NULL,
	CONSTRAINT application_profile_history_configuration_pkey PRIMARY KEY (application_profile_history_id, application_id, integration_type_id, terminal_model_id),
	CONSTRAINT application_profile_history_c_application_id_integration_t_fkey FOREIGN KEY (application_id,integration_type_id,terminal_model_id) REFERENCES application_configuration(application_id,integration_type_id,terminal_model_id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT application_profile_history_c_application_profile_history__fkey FOREIGN KEY (application_profile_history_id) REFERENCES application_profile_history(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);


-- store_go.category definition

-- Drop table

-- DROP TABLE category;

CREATE TABLE category (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	category_type_id int8 NOT NULL,
	parent_id int8 NULL,
	CONSTRAINT category_pkey PRIMARY KEY (id),
	CONSTRAINT fk_category_category_type FOREIGN KEY (category_type_id) REFERENCES category_type(id),
	CONSTRAINT fk_category_parent FOREIGN KEY (parent_id) REFERENCES category(id)
);
CREATE INDEX idx_category_ativo ON store_go.category USING btree (ativo);
CREATE INDEX idx_category_category_type_id ON store_go.category USING btree (category_type_id);
CREATE INDEX idx_category_created_at ON store_go.category USING btree (created_at);
CREATE INDEX idx_category_nome ON store_go.category USING btree (nome);
CREATE INDEX idx_category_parent_id ON store_go.category USING btree (parent_id);


-- store_go.ctgr_application_profile_history definition

-- Drop table

-- DROP TABLE ctgr_application_profile_history;

CREATE TABLE ctgr_application_profile_history (
	id_application_profile_history int8 NOT NULL,
	id_ctgr int8 NOT NULL,
	CONSTRAINT ctgr_application_profile_history_pkey PRIMARY KEY (id_application_profile_history, id_ctgr),
	CONSTRAINT fk_ctgr_application_profile_history_application_profile_history FOREIGN KEY (id_application_profile_history) REFERENCES application_profile_history(id),
	CONSTRAINT fk_ctgr_application_profile_history_category FOREIGN KEY (id_ctgr) REFERENCES category(id)
);


-- store_go.image definition

-- Drop table

-- DROP TABLE image;

CREATE TABLE image (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_obj_armazenamento int8 NOT NULL,
	cod_tip_img int2 NOT NULL,
	CONSTRAINT chk_image_cod_tip_img CHECK ((cod_tip_img = ANY (ARRAY[0, 1, 2]))),
	CONSTRAINT image_pkey PRIMARY KEY (id),
	CONSTRAINT fk_image_storage_object FOREIGN KEY (id_obj_armazenamento) REFERENCES storage_object(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_image_created_at ON store_go.image USING btree (created_at);
CREATE UNIQUE INDEX idx_image_storage_object_id ON store_go.image USING btree (id_obj_armazenamento);


-- store_go.video definition

-- Drop table

-- DROP TABLE video;

CREATE TABLE video (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_obj_armazenamento int8 NOT NULL,
	cod_tip_vid int2 NOT NULL,
	CONSTRAINT chk_video_cod_tip_vid CHECK ((cod_tip_vid = ANY (ARRAY[0, 1]))),
	CONSTRAINT video_pkey PRIMARY KEY (id),
	CONSTRAINT fk_video_storage_object FOREIGN KEY (id_obj_armazenamento) REFERENCES storage_object(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_video_created_at ON store_go.video USING btree (created_at);
CREATE UNIQUE INDEX idx_video_storage_object_id ON store_go.video USING btree (id_obj_armazenamento);


-- store_go.application_catalog definition

-- Drop table

-- DROP TABLE application_catalog;

CREATE TABLE application_catalog (
	integration_type_id int8 NOT NULL,
	terminal_model_id int8 NOT NULL,
	stage int2 NOT NULL,
	application_id int8 NOT NULL,
	ativo bool DEFAULT false NULL,
	CONSTRAINT pk_app_catalog PRIMARY KEY (integration_type_id, terminal_model_id, stage, application_id),
	CONSTRAINT fk_app_catalog_config FOREIGN KEY (application_id,integration_type_id,terminal_model_id) REFERENCES application_configuration(application_id,integration_type_id,terminal_model_id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX ak_app_ctlg_u0 ON store_go.application_catalog USING btree (application_id, stage, terminal_model_id, integration_type_id);
CREATE INDEX idx_application_catalog_ativo ON store_go.application_catalog USING btree (ativo);


-- store_go.application_version definition

-- Drop table

-- DROP TABLE application_version;

CREATE TABLE application_version (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool DEFAULT false NULL,
	application_id int8 NULL,
	integration_type_id int8 NULL,
	terminal_model_id int8 NULL,
	pilot_at timestamptz NULL,
	production_at timestamptz NULL,
	deactivated_at timestamptz NULL,
	deactivation_cause varchar(255) NULL,
	tamanho int8 NOT NULL,
	nome_versao varchar(16) NOT NULL,
	id_img int8 NULL,
	CONSTRAINT application_version_pkey PRIMARY KEY (id),
	CONSTRAINT fk_application_configuration FOREIGN KEY (application_id,integration_type_id,terminal_model_id) REFERENCES application_configuration(application_id,integration_type_id,terminal_model_id) ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT fk_application_version_image FOREIGN KEY (id_img) REFERENCES image(id)
);
CREATE INDEX ak_app_ver_u0 ON store_go.application_version USING btree (application_id, integration_type_id, terminal_model_id);
CREATE INDEX idx_application_version_ativo ON store_go.application_version USING btree (ativo);
CREATE INDEX idx_application_version_created_at ON store_go.application_version USING btree (created_at);
CREATE INDEX idx_application_version_deactivated_at ON store_go.application_version USING btree (deactivated_at);
CREATE INDEX idx_application_version_deactivation_cause ON store_go.application_version USING btree (deactivation_cause);
CREATE INDEX idx_application_version_id_image ON store_go.application_version USING btree (id_img);
CREATE INDEX idx_application_version_nome ON store_go.application_version USING btree (nome);
CREATE INDEX idx_application_version_pilot_at ON store_go.application_version USING btree (pilot_at);
CREATE INDEX idx_application_version_production_at ON store_go.application_version USING btree (production_at);