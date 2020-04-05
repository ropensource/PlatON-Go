package evm.function.paramandreturns;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.FunctionParam;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 入参是函数的使用
 * @description:
 * @author: liweic
 * @create: 2020/01/11 20:20
 **/


public class FunctionParamTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.FunctionParamTest-参数是函数的类型测试", sourcePrefix = "evm")
    public void Functionparam() {
        try {

            FunctionParam functionparam = FunctionParam.load("0xf5873e0b21ebec4e9ce615d3a9d40de8c1617b33",web3j, transactionManager, provider);

            String contractAddress = functionparam.getContractAddress();

            BigInteger t = functionparam.t().send();
            collector.logStepPass("FunctionParam函数返回值：" + t);
            collector.assertEqual("7",t.toString());

        } catch (Exception e) {
            collector.logStepFail("FunctionParamContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}



