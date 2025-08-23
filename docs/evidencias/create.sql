-- DROP SCHEMA store_go;

CREATE SCHEMA store_go AUTHORIZATION postgres;

-- DROP SEQUENCE store_go.app_cat_id_seq;

CREATE SEQUENCE store_go.app_cat_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.app_id_seq;

CREATE SEQUENCE store_go.app_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.cad_id_seq;

CREATE SEQUENCE store_go.cad_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.cat_app_id_seq;

CREATE SEQUENCE store_go.cat_app_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.cat_id_seq;

CREATE SEQUENCE store_go.cat_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.cfg_cad_id_seq;

CREATE SEQUENCE store_go.cfg_cad_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.cfg_id_seq;

CREATE SEQUENCE store_go.cfg_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.est_cat_id_seq;

CREATE SEQUENCE store_go.est_cat_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.mdl_trml_id_seq;

CREATE SEQUENCE store_go.mdl_trml_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.tip_cat_id_seq;

CREATE SEQUENCE store_go.tip_cat_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.tip_int_id_seq;

CREATE SEQUENCE store_go.tip_int_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE store_go.versao_app_id_seq;

CREATE SEQUENCE store_go.versao_app_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;-- store_go.app definition

-- Drop table

-- DROP TABLE store_go.app;

CREATE TABLE store_go.app (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	razao_social text NOT NULL,
	site text NOT NULL,
	email text NOT NULL,
	telefone text NOT NULL,
	CONSTRAINT app_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_app_created_at ON store_go.app USING btree (created_at);


-- store_go.cad definition

-- Drop table

-- DROP TABLE store_go.cad;

CREATE TABLE store_go.cad (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT cad_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_cad_created_at ON store_go.cad USING btree (created_at);


-- store_go.est_cat definition

-- Drop table

-- DROP TABLE store_go.est_cat;

CREATE TABLE store_go.est_cat (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	CONSTRAINT est_cat_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_est_cat_created_at ON store_go.est_cat USING btree (created_at);


-- store_go.mdl_trml definition

-- Drop table

-- DROP TABLE store_go.mdl_trml;

CREATE TABLE store_go.mdl_trml (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	CONSTRAINT mdl_trml_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_mdl_trml_created_at ON store_go.mdl_trml USING btree (created_at);


-- store_go.tip_cat definition

-- Drop table

-- DROP TABLE store_go.tip_cat;

CREATE TABLE store_go.tip_cat (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	CONSTRAINT tip_cat_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tip_cat_created_at ON store_go.tip_cat USING btree (created_at);


-- store_go.tip_int definition

-- Drop table

-- DROP TABLE store_go.tip_int;

CREATE TABLE store_go.tip_int (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	CONSTRAINT tip_int_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tip_int_created_at ON store_go.tip_int USING btree (created_at);


-- store_go.cat definition

-- Drop table

-- DROP TABLE store_go.cat;

CREATE TABLE store_go.cat (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	id_tipo_categoria int8 NOT NULL,
	id_pai int8 NULL,
	CONSTRAINT cat_pkey PRIMARY KEY (id),
	CONSTRAINT fk_cat_pai FOREIGN KEY (id_pai) REFERENCES store_go.cat(id),
	CONSTRAINT fk_cat_tipo_categoria FOREIGN KEY (id_tipo_categoria) REFERENCES store_go.tip_cat(id)
);
CREATE INDEX idx_cat_created_at ON store_go.cat USING btree (created_at);


-- store_go.cfg definition

-- Drop table

-- DROP TABLE store_go.cfg;

CREATE TABLE store_go.cfg (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_tipo_integracao int8 NOT NULL,
	id_modelo_terminal int8 NOT NULL,
	id_app int8 NOT NULL,
	CONSTRAINT cfg_pkey PRIMARY KEY (id),
	CONSTRAINT fk_app_configuraces FOREIGN KEY (id_app) REFERENCES store_go.app(id),
	CONSTRAINT fk_cfg_modelo_terminal FOREIGN KEY (id_modelo_terminal) REFERENCES store_go.mdl_trml(id),
	CONSTRAINT fk_cfg_tipo_integracao FOREIGN KEY (id_tipo_integracao) REFERENCES store_go.tip_int(id)
);
CREATE INDEX idx_cfg_created_at ON store_go.cfg USING btree (created_at);
CREATE UNIQUE INDEX idx_cfg_unique ON store_go.cfg USING btree (id_tipo_integracao, id_modelo_terminal, id_app);


-- store_go.cfg_cad definition

-- Drop table

-- DROP TABLE store_go.cfg_cad;

CREATE TABLE store_go.cfg_cad (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_cadastro int8 NOT NULL,
	id_configuracao int8 NOT NULL,
	CONSTRAINT cfg_cad_pkey PRIMARY KEY (id),
	CONSTRAINT fk_cad_configuracao_cadastros FOREIGN KEY (id_cadastro) REFERENCES store_go.cad(id),
	CONSTRAINT fk_cfg_configuracao_cadastros FOREIGN KEY (id_configuracao) REFERENCES store_go.cfg(id)
);
CREATE INDEX idx_cfg_cad_created_at ON store_go.cfg_cad USING btree (created_at);
CREATE INDEX idx_cfg_cad_id_configuracao ON store_go.cfg_cad USING btree (id_configuracao);


-- store_go.versao_app definition

-- Drop table

-- DROP TABLE store_go.versao_app;

CREATE TABLE store_go.versao_app (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	nome text NOT NULL,
	descricao text NULL,
	ativo bool NULL,
	id_cadastro int8 NOT NULL,
	id_configuracao int8 NOT NULL,
	tamanho int8 NOT NULL,
	nome_versao varchar(16) NOT NULL,
	CONSTRAINT versao_app_pkey PRIMARY KEY (id),
	CONSTRAINT fk_cad_versoes_aplicativo FOREIGN KEY (id_cadastro) REFERENCES store_go.cad(id),
	CONSTRAINT fk_cfg_versoes_aplicativo FOREIGN KEY (id_configuracao) REFERENCES store_go.cfg(id)
);
CREATE INDEX idx_versao_app_created_at ON store_go.versao_app USING btree (created_at);
CREATE INDEX idx_versao_app_id_configuracao ON store_go.versao_app USING btree (id_configuracao);


-- store_go.app_cat definition

-- Drop table

-- DROP TABLE store_go.app_cat;

CREATE TABLE store_go.app_cat (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_app int8 NOT NULL,
	id_categoria int8 NOT NULL,
	CONSTRAINT app_cat_pkey PRIMARY KEY (id),
	CONSTRAINT fk_app_app_categorias FOREIGN KEY (id_app) REFERENCES store_go.app(id),
	CONSTRAINT fk_app_cat_categoria FOREIGN KEY (id_categoria) REFERENCES store_go.cat(id)
);
CREATE INDEX idx_app_cat_created_at ON store_go.app_cat USING btree (created_at);
CREATE INDEX idx_app_cat_id_app ON store_go.app_cat USING btree (id_app);
CREATE INDEX idx_app_cat_id_categoria ON store_go.app_cat USING btree (id_categoria);


-- store_go.cat_app definition

-- Drop table

-- DROP TABLE store_go.cat_app;

CREATE TABLE store_go.cat_app (
	id bigserial NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	id_app int8 NOT NULL,
	id_tipo_integracao int8 NOT NULL,
	id_modelo_terminal int8 NOT NULL,
	id_estagio int8 NOT NULL,
	id_versao_aplicativo int8 NOT NULL,
	ativo bool NULL,
	CONSTRAINT cat_app_pkey PRIMARY KEY (id),
	CONSTRAINT fk_cat_app_versao_aplicativo FOREIGN KEY (id_versao_aplicativo) REFERENCES store_go.versao_app(id),
	CONSTRAINT fk_est_cat_catalogos_aplicativo FOREIGN KEY (id_estagio) REFERENCES store_go.est_cat(id)
);
CREATE INDEX idx_cat_app_created_at ON store_go.cat_app USING btree (created_at);
CREATE INDEX idx_cat_app_id_app ON store_go.cat_app USING btree (id_app);
CREATE INDEX idx_cat_app_id_estagio ON store_go.cat_app USING btree (id_estagio);
CREATE INDEX idx_cat_app_id_modelo_terminal ON store_go.cat_app USING btree (id_modelo_terminal);
CREATE INDEX idx_cat_app_id_tipo_integracao ON store_go.cat_app USING btree (id_tipo_integracao);
CREATE INDEX idx_cat_app_id_versao_aplicativo ON store_go.cat_app USING btree (id_versao_aplicativo);
CREATE INDEX idx_catapp_unique ON store_go.cat_app USING btree (id_tipo_integracao, id_modelo_terminal, id_estagio);