CREATE TABLE tmessage (
    uuid          VARCHAR(36),
    uuid_parent   VARCHAR(36),
    ishide        INTEGER DEFAULT 0,
    name	  VARCHAR(200),
    email 	  VARCHAR(200),
    userdata      VARCHAR(7000),
    text 	  VARCHAR(7000),
    data          BLOB,
    upddate       VARCHAR(20),
    tdate         VARCHAR(20),
    tdatet        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX tmessage_IDX1 ON tmessage (uuid);
CREATE        INDEX tmessage_IDX2 ON tmessage (ishide,uuid_parent,tdatet);

CREATE TABLE timage (
    uuid_message  VARCHAR(36),
    uuid          VARCHAR(36),
    hash          VARCHAR(200),
    title         VARCHAR(2000),
    path 	  VARCHAR(200),
    pathmin	  VARCHAR(200),
    imgdate       VARCHAR(40),
    imgdatet      TIMESTAMP
);
CREATE UNIQUE INDEX timage_IDX1 ON timage (uuid);
CREATE        INDEX timage_IDX2 ON timage (uuid_message);
CREATE        INDEX timage_IDX3 ON timage (hash);


COMMIT WORK;
