package evm.versioncompatible.v0_5_0.v3_deprecatedType;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DecimalLiteralsChangeByte;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 10进制数值不能直接转换成 bytesX类型
 * 必须先转换到与 bytesX相同长度的 uintY，再转换到 bytesX类型
 * @description: 
 * @author: hudenian
 * @create: 2019/12/26
 */
public class DecimalLiteralsChangeByteDeprecatedTypeTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    //需要转换成bytes4的十进制值
    private String initValue = "10";


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "version_compatible.0.5.0.constructorDeprecatedElementTest-弃用元素测试", sourcePrefix = "evm")
    public void changeByteDeprecatedType() {
        try {

            DecimalLiteralsChangeByte decimalLiteralsChangeByte = DecimalLiteralsChangeByte.load("0xc2eef35cf245edf2e87415e0b5b84951ee0a3815", web3j, transactionManager, provider);

            String contractAddress = decimalLiteralsChangeByte.getContractAddress();

            TransactionReceipt transactionReceipt = decimalLiteralsChangeByte.testChange(new BigInteger(initValue)).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            byte[] afterValueByte = decimalLiteralsChangeByte.getB4().send();

            String afterValue = DataChangeUtil.bytesToHex(afterValueByte).toLowerCase();

            collector.logStepPass(initValue+"转成bytes4后的值为："+afterValue);

            collector.assertEqual("0000000a",afterValue);
        } catch (Exception e) {
            collector.logStepFail("changeByteDeprecatedType process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
