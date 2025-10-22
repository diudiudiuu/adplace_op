import Base from './base.js';

class Client extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_client';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'client_id',
            'status',
            'expire_time',
            'api_port',
            'front_port',
            'license_key',
            'domain_deploy_flag',
            'domain_ping_flag',
            'domain_delete_flag',
            'domain_delete_days',
            'created_at',
            'updated_at',
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            client_id: {
                type: 'string',
                value: '',
            },
            status: {
                type: 'enum',
                value: ['normal', 'stop'],
                default: 'normal',
            },
            expire_time: {
                type: 'datetime',
                value: '',
                change:  [
                    'generateLicenseKey'
                ],
            },
            api_port: {
                type: 'int',
                value: 0,
            },
            front_port: {
                type: 'int',
                value: 0,
            },
            license_key: {
                type: 'string',
                value: '',
                disabled: true,
                button: {
                    text: '生成',
                    action: 'generateLicenseKey',
                },
            },
            domain_deploy_flag: {
                type: 'enum',
                value: [0, 1],
                default: 1,
            },
            domain_ping_flag: {
                type: 'enum',
                value: [0, 1],
                default: 1,
            },
            domain_delete_flag: {
                type: 'enum',
                value: [0, 1],
                default: 1,
            },
            domain_delete_days: {
                type: 'int',
                value: 14,
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

        // 最大记录数限制配置
        const maxRecords = 1; // 每个客户只能添加一个套餐
        const maxRecordsMessage = '每个客户只能添加一个套餐';

        super(tableName, primaryKey, fields, fieldTypes);
        
        // 添加限制配置
        this.maxRecords = maxRecords;
        this.maxRecordsMessage = maxRecordsMessage;
    }
}

export default Client;
