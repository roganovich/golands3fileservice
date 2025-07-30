CREATE TABLE files (
   id uuid NOT NULL, -- Идентификатор
   title varchar(1000) NULL, -- Заголовок файла
   filename varchar(1000) NOT NULL, -- Имя файла
   "extension" varchar(20) NULL, -- Расширение файла
   "size" int8 NOT NULL, -- Размер файла на диске
   date_create timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL, -- Дата создания файла
   is_deleted bool DEFAULT false NOT NULL, -- Логическое удаление
   fk_original uuid NULL, -- Ссылка на оригинал
   date_expiration date NULL, -- (DC2Type:datetimetz_immutable)
   CONSTRAINT files_pk PRIMARY KEY (id),
   CONSTRAINT original_fk FOREIGN KEY (fk_original) REFERENCES files(id)
);
CREATE INDEX idx_original_fk ON files USING btree (fk_original);
COMMENT ON TABLE files IS 'Хранение метаданных файлов';

-- Column comments

COMMENT ON COLUMN files.id IS 'Идентификатор';
COMMENT ON COLUMN files.title IS 'Заголовок файла';
COMMENT ON COLUMN files.filename IS 'Имя файла';
COMMENT ON COLUMN files."extension" IS 'Расширение файла';
COMMENT ON COLUMN files."size" IS 'Размер файла на диске';
COMMENT ON COLUMN files.date_create IS 'Дата создания файла';
COMMENT ON COLUMN files.is_deleted IS 'Логическое удаление';
COMMENT ON COLUMN files.fk_original IS 'Ссылка на оригинал';
COMMENT ON COLUMN files.date_expiration IS '(DC2Type:datetimetz_immutable)';