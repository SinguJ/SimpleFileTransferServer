import { AxiosStatic, Method } from "axios";
import { ApiResponse, ApiResult, DirMetadata, FileOrDirectory } from "../types";
// @ts-ignore
import { App } from 'vue';

class API {
  private readonly axios: AxiosStatic;
  constructor(axios: AxiosStatic) {
    this.axios = axios;
  }

  // English: Generic method for calling the API
  // 汉语：调用 API 的通用方法
  async callApi<T = ApiResult>(method: Method, uri: string) {
    // English: Format interface URI
    // 汉语：格式化接口地址
    if (!uri) {
      throw `错误: 接口地址不可为空`;
    }
    if (uri[0] !== '/') {
      uri = '/' + uri
    }
    try {
      // English: Call the API
      // 汉语：调用 API
      const resp = await this.axios.request<ApiResponse<T>>({
        url: '/api' + uri,
        method,
      });
      let data = resp.data;
      // English: Check if the call failed
      // 汉语：检查调用是否失败
      if (data.code !== 0) {
        throw `错误: 服务响应失败 => ${data.message}`;
      }
      // English: Returns the data of the interface response
      // 汉语：返回接口响应的数据
      return data.result;
    } catch (err) {
      throw err;
    }
  }

  // English: Interface to get metadata
  // 汉语：获取元数据接口
  async getMetadata(path?: string) {
    // English: Format the requested resource path
    // 汉语：格式化被请求的资源路径
    if (!path) {
      path = "/";
    } else if (path[0] !== "/") {
      path = "/" + path;
    }
    try {
      // English: Call the `get metadata interface` and return the response result
      // 汉语：调用获取元数据接口，并返回响应结果
      const res = await this.callApi<FileOrDirectory>('GET', '/api/metadata' + path);
      return res;
    } catch (err) {
      throw err;
    }
  }

  // English: An interface for reading the list of files in the specified directory
  // 汉语：读取指定目录下文件列表的接口
  async readDirectory(path?: string) {
    try {
      // English: Get the metadata of the resource corresponding to the path
      // 汉语：获取该路径对应的资源的元数据
      const res = await this.getMetadata(path);
      // English: Check if the target resource is not a directory
      // 汉语：检查目标资源是否不是目录
      if (res.type !== DirMetadata) {
        throw `错误：目标资源不是一个目录 => ${path}`;
      }
      // English: Returns a list of files in this directory
      // 汉语：返回该目录的文件列表
      return res.files;
    } catch (err) {
      throw err;
    }
  }
}

declare module "@vue/runtime-core" {
  export interface ComponentCustomProperties {
    $api: API;
  }

  export interface App {
    $api: API;
  }
}

export default function (app: App, axios: AxiosStatic): void {
  app.config.globalProperties.$api = new API(axios);
};