import Base from './base.js';

class AuthUser extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_auth_user';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'username',
            'password',
            'level',
            'status',
            'created_at',
            'updated_at',
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            username: {
                type: 'string',
                value: '',
            },
            password: {
                type: 'string',
                value: '',
            },
            level: {
                type: 'enum',
                value: ['user', 'super'],
            },
            status: {
                type: 'enum',
                value: ['normal', 'stop'],
            },
            created_at: {
                type: 'datetime',
                value: '',
            },
            updated_at: {
                type: 'datetime',
                value: '',
            }
        };

        super(tableName, primaryKey, fields, fieldTypes);
    }
}

export default AuthUser;
