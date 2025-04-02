package sqlparser_test

import (
	sqlparser "sql-parser"
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	// sql := "CREATE TABLE" + "`sometesttable` (\n" +
	// 	"`id` int(11) NOT NULL AUTO_INCREMENT ,\n" +
	// 	"`Field_1` int(11) NOT NULL DEFAULT '0' ,\n" +
	// 	"`Field_211111112` varchar(11) NOT NULL DEFAULT '' ,\n" +
	// 	"`Field_3` int(11) NOT NULL DEFAULT '0' ,\n" +
	// 	"`Field_4` int(11) NOT NULL DEFAULT '0' ,\n" +
	// 	"`Field_5` int(11) NOT NULL DEFAULT '0' ,\n" +
	// 	"`Field_6` int(11) NOT NULL DEFAULT '0' ,\n" +
	// 	"PRIMARY KEY (`id`)\n" +
	// 	") ENGINE=InnoDB DEFAULT CHARSET=utf8;\n"

	// p := sqlparser.New(strings.NewReader(sql))
	// for {
	// 	query, end := p.Parse()
	// 	if end {
	// 		break
	// 	}

	// 	t.Log(query)
	// }

	sql2 := `
CREATE TYPE session_status AS ENUM ('new', 'finished', 'active', 'declined');


CREATE TABLE sessions
(
  id        VARCHAR(255)                                     NOT NULL
    CONSTRAINT sessions_pkey
    PRIMARY KEY,
  creatorid INTEGER                                          NOT NULL,
  abonentid INTEGER                                          NOT NULL,
  status    SESSION_STATUS DEFAULT 'new' :: SESSION_STATUS   NOT NULL,
  createdat TIMESTAMP DEFAULT timezone('utc' :: TEXT, now()) NOT NULL,
  updatedat TIMESTAMP DEFAULT timezone('UTC' :: TEXT, now()) NOT NULL
);

CREATE UNIQUE INDEX sessions_id_uindex
  ON sessions (id);


CREATE OR REPLACE FUNCTION trigger_upd_time()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updatedat = (NOW() AT TIME ZONE 'UTC');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER set_upd_time
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE PROCEDURE trigger_upd_time();
	`

	p := sqlparser.New(strings.NewReader(sql2))
	for {
		query, end := p.Parse()
		if end {
			break
		}

		t.Log(query)
	}
}

func Test2(t *testing.T) {
	fs := sqlparser.NewFixedString(5)
	fs.Append('a')
	fs.Append('b')
	fs.Append('c')
	fs.Append('d')
	fs.Append('e')
	fs.Append('f')
	t.Log(fs.String())
}
