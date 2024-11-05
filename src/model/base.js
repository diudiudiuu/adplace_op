class Base {
    tableName;
    primaryKey;
    fields;
    fieldsType;
    formData;

    constructor(tableName, primaryKey, fields, fieldTypes = {}) {
        this.tableName = tableName;
        this.primaryKey = primaryKey;
        this.fields = fields;
        this.fieldsType = this.determineFieldTypes(fieldTypes);
        this.formData = {};
    }

    determineFieldTypes(fieldTypes) {
        // 默认所有字段的类型为 'string'，如果未指定则采用 fieldTypes 中的类型
        return this.fields.reduce((types, field) => {
            types[field] = fieldTypes[field] || { type: 'string', value: '' };
            return types;
        }, {});
    }

    insert() {
        const columns = this.fields.slice(1).join(', ');
        const values = this.fields.slice(1).map(field => this.escapeValue(this.formData[field]), field).join(', ');
        return `INSERT INTO ${this.tableName} (${columns}) VALUES (${values})`;
    }

    update() {
        const setClauses = this.fields.map(field => `${field} = ${this.escapeValue(this.formData[field], field)}`).join(', ');
        return `UPDATE ${this.tableName} SET ${setClauses} WHERE ${this.primaryKey} = ${this.escapeValue(this.formData[this.primaryKey], this.primaryKey)}`;
    }

    delete() {
        return `DELETE FROM ${this.tableName} WHERE ${this.primaryKey} = ${this.escapeValue(this.formData[this.primaryKey], this.primaryKey)}`;
    }

    select_list() {
        return `SELECT ${this.fields.join(', ')} FROM ${this.tableName}`;
    }

    escapeValue(value, field) {
        const fieldType = this.fieldsType[field].type;
        if (fieldType === 'string') {
            return `'${value.replace(/'/g, "''")}'`; // 转义单引号
        }
        if (fieldType === 'int') {
            return Number.parseInt(value); // 整数类型，直接返回
        }
        if (fieldType === 'enum' || fieldType === 'datetime') {
            return `'${value}'`; // 枚举和日期时间类型，返回转义后的字符串
        }
        return value; // 默认返回原始值
    }
}

export default Base;
