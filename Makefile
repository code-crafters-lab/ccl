
lint:
	@buf lint

protoc:
	@protoc -I resources \
    	--go_out . \
    	--go_opt module=github.com/code-crafters-lab/idl \
    	resources/app.proto

db-test:
	mysqldump -h 10.1.83.26 -P 3306 -u teamwork --password='jqkj5350**)' -v \
	teamwork \
	t_user t_user_bind t_role t_user_role \
	t_system t_menu t_operation t_data_limit t_api_resource t_external_link \
	t_authority t_role_authority_rel \
	t_dict t_file \
	> teamwork-test.sql

db-prod:
	mysqldump -h 192.168.44.82 -P 3306 -u teamwork --password='teamwork_jqkj5350**)123' -v \
	teamwork \
	t_user t_user_bind t_role t_user_role \
	t_system t_menu t_operation t_data_limit t_api_resource t_external_link \
	t_authority t_role_authority_rel \
	t_dict t_file \
	> teamwork-prod.sql


init:
	buf config init buf.build/ccl/dict -o dict