import { Controller, Post, Body } from '@nestjs/common';
import { AppService } from './app.service';
import { RequestDto, TransactionDto } from './serialize-params.dto';

@Controller('invokeContract')
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Post('/getWithdrawHash')
  async getWithdrawHash(@Body() request: RequestDto): Promise<string> {
    return await this.appService.getWithdrawHash(request);
  }

  @Post('/withdraw')
  async withdraw(@Body() request: RequestDto): Promise<string> {
    return await this.appService.withdraw(request);
  }

  @Post('/deposit')
  async deposit(@Body() request: RequestDto): Promise<string> {
    return await this.appService.deposit(request);
  }

  @Post('/getDepositParams')
  async getDepositParams(@Body() request: TransactionDto): Promise<object> {
    return await this.appService.getDepositParams(request);
  }

  @Post('/getBalanceOf')
  async getBalanceOf(@Body() request: RequestDto): Promise<string> {
    return await this.appService.getBalanceOf(request);
  }
}
