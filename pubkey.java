package com.chainutils.core;

import com.chainutils.crypto.Secp256k1SC;
import com.chainutils.networks.iNetwork;
import org.bouncycastle.util.BigIntegers;
import org.bouncycastle.math.ec.ECPoint;
import org.bouncycastle.util.encoders.Hex;

import java.util.Arrays;

/**
 * All rights Reserved, Designed By www.freemud.cn
 *
 * @version V1.0
 * @Title: 公钥工具类
 * @Package
 * @Description: ${TODO}(用一句话描述该文件做什么)
 * @author:
 * @date:
 * @Copyright:
 */

public class PublicKey {
    /**
     * 获取公钥
     * @param privateKey 原始私钥
     * @throws Exception
     */
    public static String getKey(byte[] privateKey, iNetwork network , boolean isCompress) {

        //生成公钥，ECDSA-secp256k1(私钥)
        ECPoint pubKey = Secp256k1SC.gMultiply(BigIntegers.fromUnsignedByteArray(privateKey));

        //获取公钥
        byte[] result = pubKey.getEncoded(isCompress);

        //如果是非压缩公钥，替换第一个byte为4
        if(!isCompress)
            result[0] = 4;

        //输出Hex格式
        return Hex.toHexString(result);
    }

    /**
     * 转压缩公钥
     * @param publicKey 未压缩公钥65位，压缩公钥33位
     * @throws Exception
     */
    public static String convertCompress(byte[] publicKey) throws Exception {
        //检查输入公钥类型
        if (isCompress(publicKey)) {
            //如果是压缩公钥直接返回
            return Hex.toHexString(publicKey);
        } else {//否则为未压缩公钥
            //定义压缩公钥33byte
            byte[] result = new byte[33];

            //获取未压缩公钥的最后一位，判断y坐标的奇偶
            int mod = Math.floorMod(publicKey[64], 2);
            //byte[0] 偶数=02，奇数=03
            if (mod == 0)
                result[0] = 2;
            else if (mod == 1)
                result[0] = 3;
            //复制x坐标，1-32byte
            for (int i = 1; i < 33; i++) {
                result[i] = publicKey[i];
            }

            return Hex.toHexString(result);
        }
    }

    /**
     * 是否为压缩公钥
     * @param publicKey 未压缩公钥65位，压缩公钥33位
     * @throws Exception
     */
    public static Boolean isCompress(byte[] publicKey) throws Exception {
        //检查输入公钥类型
        if (publicKey.length == 33) {
            //判断第一位是否是02或03
            if (publicKey[0] == 2 || publicKey[0] == 3)
                return true;
            else
                throw new Exception("输入的压缩公钥格式不正确");
        } else if (publicKey.length == 65) {
            //判断第一位是否是04
            if (publicKey[0] == 4)
                return false;
            else
                throw new Exception("输入的未压缩公钥格式不正确");
        }
        else {
            throw new Exception("输入的公钥长度不正确");
        }
    }
}
