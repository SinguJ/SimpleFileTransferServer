// English:
// <Interface Result Body>
// This interface does not have any actual properties,
// it is only used to constrain the generic declaration in the response body
//
// 汉语：
// <接口结果体>
// 该接口不存在任何实际的属性，仅用于约束响应体中的泛型声明
export interface ApiResult { }

// English: Interface Response Body
// 汉语：接口响应体
export interface ApiResponse<T = ApiResult> {
  // 状态码
  code: number;
  // 消息
  message: string;
  // 结果
  result: T;
}

// English:
// <Metadata Type>
//    0 : Indicates file metadata
//    1 : Indicates directory metadata
//
// 汉语：
// <元数据类型>
//     0 : 表示文件元数据
//     1 : 表示目录元数据
export type MetadataType = 0 | 1

// English: The metadata type value of the File
// 汉语：文件的元数据类型值
export const FileMetadata: MetadataType = 0
// English: The metadata type value of the Directory
// 汉语：目录的元数据类型值
export const DirMetadata: MetadataType = 1

// English: Metadata Interface
// 汉语：元数据接口
export interface Metadata extends ApiResult {
  // English: File name
  // 汉语：文件名称
  name: string;
  // English: File path (not the real path in the file system,
  //          only relative to the file root path of the web service)
  // 汉语：文件路径（非文件系统中的真实路径，仅相对于 Web 服务的文件根路径）
  path: string;
  // English: Metadata Type Identifier
  // 汉语：元数据类型标识
  type: MetadataType;
}

// English: Directory Metadata Interface
// 汉语：目录元数据接口
export interface Directory extends Metadata {
  // English: Child metadata list
  // 汉语：子元数据列表
  files: Array<File | Directory>;
}

// English: File Metadata Interface
// 汉语：文件元数据接口
export interface File extends Metadata {
  // English: File size (unit: byte)
  // 汉语：文件大小（单位：字节）
  size: number;
  // English: File type
  // 汉语：文件类型
  filetype: string;
}


// English: 
// <File Or Directory Metadata Interface>
// In order to eliminate frequent type conversions,
// this type of fusion of File and Directory properties is defined
//
// 汉语：
// <文件或目录元数据接口>
// 为了消除频繁的类型转换，定义了这种融合了 File 和 Directory 属性的类型
export interface FileOrDirectory extends Metadata {
  // English: Child metadata list
  // 汉语：子元数据列表
  files?: Array<FileOrDirectory>;
  // English: File size (unit: byte)
  // 汉语：文件大小（单位：字节）
  size?: number;
  // English: File type
  // 汉语：文件类型
  filetype?: string;
}