import Base from './base.js';

class AuthUsertoken extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_auth_user_token';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'auth_user_id',
            'token',
            'expires_at',
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            auth_user_id: {
                type: 'int',
                value: '',
            },
            token: {
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
            expires_at: {
                type: 'datetime',
                value: '',
            }
        };

        super(tableName, primaryKey, fields, fieldTypes);
    }
}

export default AuthUsertoken;
