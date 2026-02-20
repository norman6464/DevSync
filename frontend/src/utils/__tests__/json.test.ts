import { describe, it, expect } from 'vitest';
import { parseJsonArray, parseJsonObject } from '../json';

describe('parseJsonArray', () => {
  it('正常なJSON配列文字列をパースする', () => {
    expect(parseJsonArray('["a","b","c"]')).toEqual(['a', 'b', 'c']);
  });

  it('数値配列をパースする', () => {
    expect(parseJsonArray<number>('[1,2,3]')).toEqual([1, 2, 3]);
  });

  it('空配列のJSON文字列を正しく処理する', () => {
    expect(parseJsonArray('[]')).toEqual([]);
  });

  it('undefinedの場合に空配列を返す', () => {
    expect(parseJsonArray(undefined)).toEqual([]);
  });

  it('nullの場合に空配列を返す', () => {
    expect(parseJsonArray(null)).toEqual([]);
  });

  it('空文字列の場合に空配列を返す', () => {
    expect(parseJsonArray('')).toEqual([]);
  });

  it('不正なJSONの場合に空配列を返す', () => {
    expect(parseJsonArray('invalid json')).toEqual([]);
  });

  it('途中で切れたJSONの場合に空配列を返す', () => {
    expect(parseJsonArray('["a","b"')).toEqual([]);
  });

  it('ネストされたオブジェクト配列をパースする', () => {
    const input = '[{"id":1,"name":"test"}]';
    expect(parseJsonArray<{ id: number; name: string }>(input)).toEqual([{ id: 1, name: 'test' }]);
  });

  it('日本語を含む配列をパースする', () => {
    expect(parseJsonArray('["タグ1","タグ2"]')).toEqual(['タグ1', 'タグ2']);
  });
});

describe('parseJsonObject', () => {
  it('正常なJSONオブジェクトをパースする', () => {
    expect(parseJsonObject('{"key":"value"}')).toEqual({ key: 'value' });
  });

  it('ネストされたオブジェクトをパースする', () => {
    expect(parseJsonObject('{"a":{"b":1}}')).toEqual({ a: { b: 1 } });
  });

  it('undefinedの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject(undefined)).toEqual({});
  });

  it('nullの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject(null)).toEqual({});
  });

  it('空文字列の場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('')).toEqual({});
  });

  it('不正なJSONの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('invalid json')).toEqual({});
  });

  it('JSON配列の場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('[1,2,3]')).toEqual({});
  });

  it('JSON文字列リテラルの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('"hello"')).toEqual({});
  });

  it('JSON nullの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('null')).toEqual({});
  });

  it('数値のみの場合に空オブジェクトを返す', () => {
    expect(parseJsonObject('42')).toEqual({});
  });
});
